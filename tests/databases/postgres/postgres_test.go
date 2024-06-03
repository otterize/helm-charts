package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	tests "helm_tests/databases"
	"io"
	"k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
	"testing"
	"time"
)

const (
	OtterizeKubernetesChartPath = "../../../otterize-kubernetes"
	TestNamespace               = "postgres-integration-test"
	PostgresRootPassword        = "integrationtestpassword11"
	PostgresCredsSecretName     = "postgres-user-password"
	PostgresSvcName             = "otterize-database"
	PostgresDatabaseName        = "test-db"
	PostgresInstanceName        = "otterize-postgres"
	PostgresRootUser            = "otterize-admin"
	IntentsResourceName         = "psql-client-intents"
	PostgresConnectionString    = "postgres://%s:%s@%s:5432/%s"
)

var OtterizeValuesMapperDisabled = map[string]interface{}{
	"global": map[string]interface{}{
		"deployment": map[string]interface{}{
			"networkMapper": false,
		},
	},
}

type PostgresTestSuite struct {
	tests.BaseSuite
	clientPodName string
}

func (s *PostgresTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	logrus.Info("Setting up postgres test suite")
	logrus.Info("Installing otterize-kubernetes helm chart")
	// Load Chart.yaml
	chart, err := loader.Load(OtterizeKubernetesChartPath)
	s.Require().NoError(err)
	//
	installAction := action.NewInstall(s.HelmActionConfig)
	installAction.Namespace = "otterize-system"
	installAction.ReleaseName = "otterize"
	installAction.CreateNamespace = true
	installAction.Wait = true
	installAction.Timeout = time.Second * 40

	// Run helm install command
	_, err = installAction.Run(chart, OtterizeValuesMapperDisabled)

	s.Require().NoError(err)
	logrus.Info("otterize-kubernetes started successfully")

	// Create test namespace
	logrus.Info("Creating test namespace")
	_, err = s.Client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: TestNamespace},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	// Deploy postgres pod & service
	logrus.Info("Deploying PostgreSQL database pod & service")
	s.deployAndConfigureDatabase()

	// Deploy client pod
	logrus.Info("Deploying psql client pod")
	s.deployDatabaseClient()
	s.waitForPodToStart("app=psql-client", 30) // Dependent on a secret to be created

	// Get client pod name
	res, err := s.Client.CoreV1().Pods(TestNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: "app=psql-client"})
	s.Require().NoError(err)
	s.Require().Len(res.Items, 1)
	s.clientPodName = res.Items[0].Name

	s.applyPGServerConf()

	logrus.Info("PostgreSQL Suite setup complete")
}

func (s *PostgresTestSuite) SetupTest() {
	err := s.IntentsClient.Namespace(TestNamespace).Delete(context.Background(), IntentsResourceName, metav1.DeleteOptions{})
	if !errors.IsNotFound(err) {
		s.Require().NoError(err) // Just fail
	}
}

func (s *PostgresTestSuite) TestWorkloadFailsToAccessDatabase() {
	logrus.Info("Validating client pod fails to access the database")
	s.matchSubStringsInLog(s.clientPodName, TestNamespace, []string{"password authentication failed"})
}

func (s *PostgresTestSuite) TestAddSelectAndInsertPermissionsForDB() {
	s.applyIntents([]v1alpha3.DatabaseOperation{v1alpha3.DatabaseOperationInsert, v1alpha3.DatabaseOperationSelect})
	logrus.Info("Validating client pod was granted SELECT & INSERT permissions")
	s.matchSubStringsInLog(s.clientPodName, TestNamespace, []string{"Successfully INSERTED", "Successfully SELECTED"})
}

func (s *PostgresTestSuite) TearDownSuite() {
	err := s.IntentsClient.Namespace(TestNamespace).Delete(context.Background(), IntentsResourceName, metav1.DeleteOptions{})
	if !errors.IsNotFound(err) {
		s.Require().NoError(err) // Just fail
	}

	err = s.Client.CoreV1().Namespaces().Delete(context.Background(), TestNamespace, metav1.DeleteOptions{})
	if !errors.IsNotFound(err) {
		s.Require().NoError(err) // Just fail
	}

	uninstallAction := action.NewUninstall(s.HelmActionConfig)
	_, err = uninstallAction.Run("otterize")
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) deployAndConfigureDatabase() {
	postgresDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PostgresInstanceName,
			Namespace: TestNamespace,
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
	_, err := s.Client.AppsV1().Deployments(TestNamespace).Create(context.Background(), postgresDeployment, metav1.CreateOptions{})
	s.Require().NoError(err)

	postgresService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PostgresSvcName,
			Namespace: TestNamespace,
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
	_, err = s.Client.CoreV1().Services(TestNamespace).Create(context.Background(), &postgresService, metav1.CreateOptions{})
	s.Require().NoError(err)

	s.waitForPodToStart("app=database", 10)

	logrus.Info("Spawning job to create a test table in the database")
	s.runCreateTableJob()
}

func (s *PostgresTestSuite) applyIntents(operations []v1alpha3.DatabaseOperation) {
	clientIntents := v1alpha3.ClientIntents{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClientIntents",
			APIVersion: "k8s.otterize.com/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      IntentsResourceName,
			Namespace: TestNamespace,
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
	u := s.getUnstructuredObject(clientIntents)
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "k8s.otterize.com",
		Version: "v1alpha3",
		Kind:    "ClientIntents",
	})
	_, err := s.IntentsClient.Namespace(TestNamespace).Create(context.Background(), u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) deployDatabaseClient() {
	clientDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "psql-client",
			Namespace: TestNamespace,
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
	_, err := s.Client.AppsV1().Deployments(TestNamespace).Create(context.Background(), clientDeployment, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) matchSubStringsInLog(podName, podNamespace string, stringsToMatch []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// The client is sampling DB access every 5 seconds, we do the same for the logs
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			logrus.Infof("Trying to match %v in pod logs", stringsToMatch)
			logMessage := s.getPodLog(podName, podNamespace)
			if logMessage != "" {
				var matched = true
				for _, stringToMatch := range stringsToMatch {
					if !strings.Contains(logMessage, stringToMatch) {
						matched = false
					}
				}
				if matched {
					return
				}
			}
		case <-ctx.Done():
			s.FailNowf("Could not match all strings in log", "Failed matching: %v in psql-client logs", stringsToMatch)
		}
	}
}

func (s *PostgresTestSuite) waitForPodToStart(podLabelSelector string, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*duration)
	defer cancel()
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			res, err := s.Client.CoreV1().Pods(TestNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: podLabelSelector})
			s.Require().NoError(err)
			if len(res.Items) == 0 {
				continue
			}
			pod := res.Items[0]
			logrus.Infof("%s is in phase %s", pod.Name, pod.Status.Phase)
			if pod.Status.Phase == corev1.PodRunning {
				return
			}

		case <-ctx.Done():
			s.FailNowf("Pods in namespace were not ready", "Pods were not running after %d seconds", duration)
		}
	}
}

func (s *PostgresTestSuite) applyPGServerConf() {
	pgServerConf := v1alpha3.PostgreSQLServerConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PostgreSQLServerConfig",
			APIVersion: "k8s.otterize.com/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: PostgresInstanceName,
		},
		Spec: v1alpha3.PostgreSQLServerConfigSpec{
			Address: fmt.Sprintf("%s.%s.svc.cluster.local:5432", PostgresSvcName, TestNamespace),
			Credentials: v1alpha3.DatabaseCredentials{
				Username: PostgresRootUser,
				Password: PostgresRootPassword,
			},
		},
	}

	u := s.getUnstructuredObject(pgServerConf)
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "k8s.otterize.com",
		Version: "v1alpha3",
		Kind:    "PostgreSQLServerConfig",
	})

	_, err := s.PGServerConfClient.Namespace(TestNamespace).Create(context.Background(), u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) getUnstructuredObject(resource any) *unstructured.Unstructured {
	body, err := json.Marshal(resource)
	s.Require().NoError(err)

	u := unstructured.Unstructured{}
	err = u.UnmarshalJSON(body)
	s.Require().NoError(err)

	return &u
}

func (s *PostgresTestSuite) runCreateTableJob() {
	connectionString := fmt.Sprintf(PostgresConnectionString, PostgresRootUser, PostgresRootPassword, PostgresSvcName, PostgresDatabaseName)
	res, err := s.Client.BatchV1().Jobs(TestNamespace).Create(context.Background(), &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "create-table-job",
			Namespace: TestNamespace,
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

func (s *PostgresTestSuite) getPodLog(name string, namespace string) string {
	linesToTail := int64(10)
	req := s.Client.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		Follow:    true,
		TailLines: &linesToTail,
	})
	podLogs, err := req.Stream(context.Background())
	s.Require().NoError(err)
	defer podLogs.Close()
	buf := make([]byte, 200)
	numBytes, err := podLogs.Read(buf)
	s.Require().NoError(err)
	if numBytes == 0 {
		return ""
	}
	if err == io.EOF {
		return ""
	}
	return string(buf[:numBytes])
}

func TestPostgresEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
