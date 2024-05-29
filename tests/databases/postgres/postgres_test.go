package postgres

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	tests "helm_tests"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

const (
	OtterizeKubernetesChartPath = "../../otterize-kubernetes"
	TestNamespace               = "postgres-integration-test"
	PostgresRootPassword        = "integrationtestpassword11"
	PostgresCredsSecretName     = "postgres-user-password"
)

type PostgresTestSuite struct {
	tests.BaseSuite
}

func (s *PostgresTestSuite) SetupTest() {
	fmt.Println("SetupTest")
}

func (s *PostgresTestSuite) TestOtterizeKubernetesHelmInstall() {
	// Load Chart.yaml
	chart, err := loader.Load(OtterizeKubernetesChartPath)
	s.Require().NoError(err)
	//
	installAction := action.NewInstall(s.HelmActionConfig)
	installAction.Namespace = "otterize-system"
	installAction.ReleaseName = "otterize"
	installAction.CreateNamespace = true
	installAction.Wait = true

	// Run helm install command
	_, err = installAction.Run(chart, chart.Values)
	s.Require().NoError(err)

	// Create test namespace
	_, err = s.Client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: TestNamespace},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	// Deploy postgres pod & service
	s.deployPostgresDatabase()

	// Deploy client pod
	s.deployDatabaseClient()

}

func (s *PostgresTestSuite) TearDownTest() {
	uninstallAction := action.NewUninstall(s.HelmActionConfig)
	_, err := uninstallAction.Run("otterize")
	s.Require().NoError(err)

	err = s.Client.CoreV1().Namespaces().Delete(context.Background(), TestNamespace, metav1.DeleteOptions{})
	s.Require().NoError(err)

}

func (s *PostgresTestSuite) deployPostgresDatabase() {
	// Deployment
	postgresDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
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
								Value: "otterize-test",
							},
							{
								Name:  "POSTGRES_USER",
								Value: "otterize-admin",
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

	// Service
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "otterize-database",
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
	_, err = s.Client.CoreV1().Services(TestNamespace).Create(context.Background(), &service, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) applyCRD() {
	//
	//obj := &unstructured.Unstructured{
	//	Object: map[string]interface{}{
	//		"apiVersion": "flink.apache.org/v1beta1",
	//		"kind":       "FlinkSessionJob",
	//		"metadata": map[string]interface{}{
	//			"name": "name of your application",
	//		},
	//		"spec": map[string]interface{}{
	//			"deploymentName": "Deployment name to which this FlinkSessionJob belongs",
	//			"job": map[string]interface{}{
	//				"jarURI":      "Your Jar file",
	//				"parallelism": 4,
	//				"upgradeMode": "upgrade mode",
	//				"state":       "running",
	//			},
	//		},
	//	},
	//}
}

func (s *PostgresTestSuite) deployDatabaseClient() {
	postgresDeployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "psql-client",
			Namespace: TestNamespace,
			Annotations: map[string]string{
				"credentials-operator.otterize.com/user-password-secret-name": PostgresCredsSecretName,
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "psql-client",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "psql-client"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "client",
						Image: "otterize/postgres-client",
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
								Value: "otterize-database",
							},
							{
								Name:  "DATABASE_NAME",
								Value: "otterize-test",
							},
							{
								Name:  "DATABASE_PORT",
								Value: "5432",
							},
						},
					}},
				},
			},
		},
	}
	_, err := s.Client.AppsV1().Deployments(TestNamespace).Create(context.Background(), postgresDeployment, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func TestPostgresEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
