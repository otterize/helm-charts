package postgres

import (
	"context"
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm_tests"
	"k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"testing"
	"time"
)

const (
	PostgresRootPassword     = "integrationtestpassword11"
	PostgresCredsSecretName  = "postgres-user-password"
	PostgresSvcName          = "otterize-database"
	PostgresDatabaseName     = "test-db"
	PostgresInstanceName     = "otterize-postgres"
	PostgresRootUser         = "otterize-admin"
	IntentsResourceName      = "psql-client-intents"
	PostgresConnectionString = "postgres://%s:%s@%s:5432/%s"
)

type PostgresTestSuite struct {
	helm_tests.BaseSuite
	PGServerConfClient dynamic.NamespaceableResourceInterface
	clientPod          *corev1.Pod
	testNamespaceName  string
}

func (s *PostgresTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.InstallOtterizeHelmChart(s.GetDefaultHelmChartValues())

	s.PGServerConfClient = s.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "k8s.otterize.com",
		Version:  "v1alpha3",
		Resource: "postgresqlserverconfigs",
	})
}

func (s *PostgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()
	s.UninstallOtterizeHelmChart(ctx)
}

func (s *PostgresTestSuite) SetupTest() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()

	// Create test namespace
	s.testNamespaceName = fmt.Sprintf("postgres-integration-test-%d", time.Now().Unix())
	logrus.Info("Creating test namespace")
	s.CreateNamespace(ctx, s.testNamespaceName)

	// Deploy postgres pod & service
	logrus.Info("Deploying PostgreSQL database pod & service")
	s.deployAndConfigureDatabase(ctx)

	// Deploy client pod
	logrus.Info("Deploying psql client pod")
	s.deployDatabaseClient(ctx)

	// Get client pod name
	s.clientPod = s.FindPodByLabel(ctx, s.testNamespaceName, "app=psql-client")

	s.applyPGServerConf(ctx)
}

func (s *PostgresTestSuite) TearDownTest() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()
	// ClientIntents have to be deleted before the namespace as properly deleting them requires the existence of the PGServerConf
	s.DeleteClientIntents(ctx, s.testNamespaceName, IntentsResourceName)
	s.DeleteNamespace(ctx, s.testNamespaceName)
}

func (s *PostgresTestSuite) deployPostgresDatabase(ctx context.Context) {
	postgresDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PostgresInstanceName,
			Namespace: s.testNamespaceName,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "database",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "database"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "database",
						Image: "postgres:latest",
						Env: []corev1.EnvVar{
							{
								Name:  "POSTGRES_DB",
								Value: PostgresDatabaseName,
							},
							{
								Name:  "POSTGRES_USER",
								Value: PostgresRootUser,
							},
							{
								Name:  "POSTGRES_PASSWORD",
								Value: PostgresRootPassword,
							},
						},
					}},
				},
			},
		},
	}
	s.CreateDeployment(ctx, postgresDeployment)
	s.WaitForDeploymentAvailability(ctx, postgresDeployment.Namespace, postgresDeployment.Name)
}

func (s *PostgresTestSuite) createPostgresService(ctx context.Context) {
	postgresService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PostgresSvcName,
			Namespace: s.testNamespaceName,
			Labels: map[string]string{
				"app": "database-svc",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       5432,
					TargetPort: intstr.IntOrString{IntVal: 5432},
				},
			},
			Selector: map[string]string{"app": "database"},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	s.CreateService(ctx, &postgresService)
}

func (s *PostgresTestSuite) deployAndConfigureDatabase(ctx context.Context) {
	s.deployPostgresDatabase(ctx)
	s.createPostgresService(ctx)

	logrus.Info("Spawning job to create a test table in the database")
	s.runCreateTableJob(ctx)
}

func (s *PostgresTestSuite) applyIntents(ctx context.Context, operations []v1alpha3.DatabaseOperation) {
	clientIntents := v1alpha3.ClientIntents{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClientIntents",
			APIVersion: "k8s.otterize.com/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      IntentsResourceName,
			Namespace: s.testNamespaceName,
		},
		Spec: &v1alpha3.IntentsSpec{
			Service: v1alpha3.Service{
				Name: "psql-client",
			},
			Calls: []v1alpha3.Intent{
				{
					Name: PostgresInstanceName,
					Type: v1alpha3.IntentTypeDatabase,
					DatabaseResources: []v1alpha3.DatabaseResource{
						{
							DatabaseName: PostgresDatabaseName,
							Operations:   operations,
						},
					},
				},
			},
		},
		Status: v1alpha3.IntentsStatus{},
	}

	s.ApplyClientIntents(ctx, clientIntents)
}

func (s *PostgresTestSuite) deployDatabaseClient(ctx context.Context) {
	clientDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "psql-client",
			Namespace: s.testNamespaceName,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "psql-client",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"credentials-operator.otterize.com/user-password-secret-name": PostgresCredsSecretName,
					},
					Labels: map[string]string{"app": "psql-client"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "client",
						Image: "otterize/postgres-integration-test-client",
						Env: []corev1.EnvVar{
							{
								Name: "DATABASE_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: PostgresCredsSecretName,
										},
										Key: "username",
									},
								},
							},
							{
								Name: "DATABASE_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: PostgresCredsSecretName,
										},
										Key: "password",
									},
								},
							},
							{
								Name:  "DATABASE_HOST",
								Value: PostgresSvcName,
							},
							{
								Name:  "DATABASE_NAME",
								Value: PostgresDatabaseName,
							},
						},
					}},
				},
			},
		},
	}

	s.CreateDeployment(ctx, clientDeployment)
	s.WaitForDeploymentAvailability(ctx, clientDeployment.Namespace, clientDeployment.Name)
}

func (s *PostgresTestSuite) applyPGServerConf(ctx context.Context) {
	pgServerConf := v1alpha3.PostgreSQLServerConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PostgreSQLServerConfig",
			APIVersion: "k8s.otterize.com/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: PostgresInstanceName,
		},
		Spec: v1alpha3.PostgreSQLServerConfigSpec{
			Address: fmt.Sprintf("%s.%s.svc.cluster.local:5432", PostgresSvcName, s.testNamespaceName),
			Credentials: v1alpha3.DatabaseCredentials{
				Username: PostgresRootUser,
				Password: PostgresRootPassword,
			},
		},
	}

	u := s.GetUnstructuredObject(pgServerConf, pgServerConf.GroupVersionKind())
	_, err := s.PGServerConfClient.Namespace(s.testNamespaceName).Create(ctx, u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) runCreateTableJob(ctx context.Context) {
	connectionString := fmt.Sprintf(PostgresConnectionString, PostgresRootUser, PostgresRootPassword, PostgresSvcName, PostgresDatabaseName)
	res, err := s.Client.BatchV1().Jobs(s.testNamespaceName).Create(ctx, &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "create-table-job",
			Namespace: s.testNamespaceName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:    "create-table",
							Image:   "postgres:latest",
							Command: []string{"psql"},
							Args:    []string{connectionString, "-c", "CREATE TABLE IF NOT EXISTS example ( entry_time BIGINT );"},
						},
					},
				},
			},
		},
		Status: batchv1.JobStatus{},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
}

func (s *PostgresTestSuite) TestWorkloadFailsToAccessDatabase() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	logrus.Info("Validating client pod fails to access the database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "password authentication failed")
}

func (s *PostgresTestSuite) TestAddSelectAndInsertPermissionsForDB() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "password authentication failed")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert, v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT & INSERT permissions")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully INSERTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully SELECTED")
}

func (s *PostgresTestSuite) TestInsertPermissionWithoutSelect() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "password authentication failed")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert})
	logrus.Info("Validating client pod was granted INSERT permissions without SELECT")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully INSERTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Unable to perform SELECT operation")
}

func (s *PostgresTestSuite) TestSelectPermissionWithoutInsert() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "password authentication failed")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT permissions without INSERT")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully SELECTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Unable to perform INSERT operation")
}

func TestPostgresEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
