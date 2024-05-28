package databases

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
	"testing"
)

const (
	OtterizeKubernetesChartPath = "../../otterize-kubernetes"
	TestNamespace               = "postgres-integration-test"
	PostgresRootPassword        = "integrationtestpassword11"
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
	// Run helm install command
	_, err = installAction.Run(chart, chart.Values)
	s.Require().NoError(err)

	// Deploy database pod
	s.deployPostgreSQLDatabasePod()

	// Deploy client pod

}

func (s *PostgresTestSuite) TearDownTest() {
	err := s.Client.CoreV1().Namespaces().Delete(context.Background(), TestNamespace, metav1.DeleteOptions{})
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) deployPostgreSQLDatabasePod() {
	_, err := s.Client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: TestNamespace},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	databasePod := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database",
			Namespace: TestNamespace,
		},
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
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
	_, err = s.Client.AppsV1().Deployments(TestNamespace).Create(context.Background(), databasePod, metav1.CreateOptions{})
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

func TestPostgresEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
