package mysql

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
	MySQLRootCredentialsSecretName = "mysql-root-credentials"
	MySQLRootPassword              = "password"
	MySQLRootUser                  = "root"
	MySQLCredsSecretName           = "mysql-user-password"
	MySQLSvcName                   = "otterize-database"
	MySQLDatabaseName              = "testdb"
	MySQLInstanceName              = "otterize-mysql"
	IntentsResourceName            = "mysql-client-intents"
	MySQLImage                     = "mysql:8.0"
)

type MySQLTestSuite struct {
	helm_tests.BaseSuite
	MySQLServerConfClient dynamic.NamespaceableResourceInterface
	clientPod             *corev1.Pod
	testNamespaceName     string
}

func (s *MySQLTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.InstallOtterizeHelmChart(s.GetDefaultHelmChartValues())

	s.MySQLServerConfClient = s.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "k8s.otterize.com",
		Version:  "v1alpha3",
		Resource: "mysqlserverconfigs",
	})
}

func (s *MySQLTestSuite) TearDownSuite() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()
	s.UninstallOtterizeHelmChart(ctx)
}

func (s *MySQLTestSuite) SetupTest() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()

	// Create test namespace
	s.testNamespaceName = fmt.Sprintf("mysql-integration-test-%d", time.Now().Unix())
	logrus.Info("Creating test namespace")
	s.CreateNamespace(ctx, s.testNamespaceName)

	// Deploy mysql pod & service
	logrus.Info("Deploying MySQL database pod & service")
	s.deployAndConfigureDatabase(ctx)

	// Deploy client pod
	logrus.Info("Deploying mysql client pod")
	s.deployDatabaseClient(ctx)

	// Get client pod name
	s.clientPod = s.FindPodByLabel(ctx, s.testNamespaceName, "app=mysql-client")
}

func (s *MySQLTestSuite) TearDownTest() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	defer cancel()
	// ClientIntents have to be deleted before the namespace as properly deleting them requires the existence of the MySQLServerConf
	s.DeleteClientIntents(ctx, s.testNamespaceName, IntentsResourceName)
	s.DeleteNamespace(ctx, s.testNamespaceName)
}

func (s *MySQLTestSuite) deployMySQLDatabase(ctx context.Context) {
	mysqlDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MySQLInstanceName,
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
						Image: MySQLImage,
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: MySQLRootPassword,
							},
							{
								Name:  "MYSQL_ROOT_HOST",
								Value: "%",
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: MySQLDatabaseName,
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},

						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"mysqladmin", "ping", "-h", "localhost", "--protocol=tcp"},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       1,
							TimeoutSeconds:      1,
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"mysqladmin", "ping", "-h", "localhost", "--protocol=tcp"},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       1,
							TimeoutSeconds:      1,
						},
					}},
				},
			},
		},
	}
	s.CreateDeployment(ctx, mysqlDeployment)
	s.WaitForDeploymentAvailability(ctx, mysqlDeployment.Namespace, mysqlDeployment.Name)
}

func (s *MySQLTestSuite) createMySQLService(ctx context.Context) {
	mySQLService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MySQLSvcName,
			Namespace: s.testNamespaceName,
			Labels: map[string]string{
				"app": "database-svc",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					TargetPort: intstr.IntOrString{IntVal: 3306},
				},
			},
			Selector: map[string]string{"app": "database"},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	s.CreateService(ctx, &mySQLService)
}

func (s *MySQLTestSuite) deployAndConfigureDatabase(ctx context.Context) {
	s.deployMySQLDatabase(ctx)
	s.createMySQLService(ctx)

	logrus.Info("Spawning job to create a test table in the database")
	s.runCreateTableJob(ctx)
}

func (s *MySQLTestSuite) applyIntents(ctx context.Context, operations []v1alpha3.DatabaseOperation) {
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
				Name: "mysql-client",
			},
			Calls: []v1alpha3.Intent{
				{
					Name: MySQLInstanceName,
					Type: v1alpha3.IntentTypeDatabase,
					DatabaseResources: []v1alpha3.DatabaseResource{
						{
							DatabaseName: MySQLDatabaseName,
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

func (s *MySQLTestSuite) deployDatabaseClient(ctx context.Context) {
	clientDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-client",
			Namespace: s.testNamespaceName,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mysql-client",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"credentials-operator.otterize.com/user-password-secret-name": MySQLCredsSecretName,
					},
					Labels: map[string]string{"app": "mysql-client"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "client",
						Image: "otterize/mysql-integration-test-client:latest",
						Env: []corev1.EnvVar{
							{
								Name: "DATABASE_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: MySQLCredsSecretName,
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
											Name: MySQLCredsSecretName,
										},
										Key: "password",
									},
								},
							},
							{
								Name:  "DATABASE_HOST",
								Value: MySQLSvcName,
							},
							{
								Name:  "DATABASE_NAME",
								Value: MySQLDatabaseName,
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

func (s *MySQLTestSuite) CreateMySQLServerConf(ctx context.Context, mySQLServerConf *v1alpha3.MySQLServerConfig) {
	logrus.WithField("namespace", mySQLServerConf.Namespace).WithField("name", mySQLServerConf.Name).Info("Creating MySQLServerConfig")
	mySQLServerConf.TypeMeta = metav1.TypeMeta{
		Kind:       "MySQLServerConfig",
		APIVersion: "k8s.otterize.com/v1alpha3",
	}
	u := s.GetUnstructuredObject(mySQLServerConf, mySQLServerConf.GroupVersionKind())
	_, err := s.MySQLServerConfClient.Namespace(mySQLServerConf.Namespace).Create(ctx, u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *MySQLTestSuite) applyMySQLServerConfWithInlinePassword(ctx context.Context) {
	mySQLServerConf := v1alpha3.MySQLServerConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MySQLInstanceName,
			Namespace: s.testNamespaceName,
		},
		Spec: v1alpha3.MySQLServerConfigSpec{
			Address: fmt.Sprintf("%s.%s.svc.cluster.local:3306", MySQLSvcName, s.testNamespaceName),
			Credentials: v1alpha3.DatabaseCredentials{
				Username: MySQLRootUser,
				Password: MySQLRootPassword,
			},
		},
	}

	s.CreateMySQLServerConf(ctx, &mySQLServerConf)
}

func (s *MySQLTestSuite) applyMySQLServerConfWithSecretRef(ctx context.Context) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MySQLRootCredentialsSecretName,
			Namespace: s.testNamespaceName,
		},
		StringData: map[string]string{
			"username": MySQLRootUser,
			"password": MySQLRootPassword,
		},
	}
	s.CreateSecret(ctx, secret)

	mySQLServerConf := v1alpha3.MySQLServerConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MySQLInstanceName,
			Namespace: s.testNamespaceName,
		},
		Spec: v1alpha3.MySQLServerConfigSpec{
			Address: fmt.Sprintf("%s.%s.svc.cluster.local:3306", MySQLSvcName, s.testNamespaceName),
			Credentials: v1alpha3.DatabaseCredentials{
				SecretRef: &v1alpha3.DatabaseCredentialsSecretRef{
					Name: MySQLRootCredentialsSecretName,
				},
			},
		},
	}
	s.CreateMySQLServerConf(ctx, &mySQLServerConf)
}

func (s *MySQLTestSuite) runCreateTableJob(ctx context.Context) {
	job := &batchv1.Job{
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
							Image:   MySQLImage,
							Command: []string{"mysql"},
							Args:    []string{"-h", MySQLSvcName, "-u", MySQLRootUser, "--database", MySQLDatabaseName, "-e", "CREATE TABLE IF NOT EXISTS example ( entry_time BIGINT );"},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_PWD",
									Value: MySQLRootPassword,
								},
							},
						},
					},
				},
			},
		},
		Status: batchv1.JobStatus{},
	}
	s.CreateJob(ctx, job)
	s.WaitForJobCompletion(ctx, job.Namespace, job.Name)
}

func (s *MySQLTestSuite) TestWorkloadFailsToAccessDatabase() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.applyMySQLServerConfWithInlinePassword(ctx)

	logrus.Info("Validating client pod fails to access the database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Access denied for user")
}

func (s *MySQLTestSuite) TestAddSelectAndInsertPermissionsForDB() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.applyMySQLServerConfWithInlinePassword(ctx)

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Access denied for user")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert, v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT & INSERT permissions")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully INSERTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully SELECTED")
}

func (s *MySQLTestSuite) TestInsertPermissionWithoutSelect() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.applyMySQLServerConfWithInlinePassword(ctx)

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Access denied for user")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert})
	logrus.Info("Validating client pod was granted INSERT permissions without SELECT")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully INSERTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Unable to perform SELECT operation")
}

func (s *MySQLTestSuite) TestSelectPermissionWithoutInsert() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.applyMySQLServerConfWithInlinePassword(ctx)

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Access denied for user")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT permissions without INSERT")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully SELECTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Unable to perform INSERT operation")
}

func (s *MySQLTestSuite) TestAddSelectAndInsertPermissionsWithSecretRefPermissions() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	s.applyMySQLServerConfWithSecretRef(ctx)

	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Access denied for user")

	s.applyIntents(ctx, []v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert, v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT & INSERT permissions")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully connected to database")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully INSERTED")
	s.ReadPodLogsUntilSubstring(ctx, s.clientPod, "Successfully SELECTED")
}

func TestMySQLEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(MySQLTestSuite))
}
