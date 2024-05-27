package databases

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm_tests/config"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"testing"
)

const (
	OtterizeKubernetesChartPath = "../../otterize-kubernetes"
	TestNamespace               = "postgres-integration-test"
	PostgresRootPassword        = "integrationtestpassword11"
)

type PostgresTestSuite struct {
	suite.Suite
	client           *kubernetes.Clientset
	helmActionConfig *action.Configuration
}

func (s *PostgresTestSuite) SetupTest() {
	homeDir, err := os.UserHomeDir()
	s.Require().NoError(err)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/%s", homeDir, viper.GetString(config.KubeConfigPath)))
	s.Require().NoError(err)
	//
	client, err := kubernetes.NewForConfig(kubeConfig)
	s.Require().NoError(err)
	s.client = client

	actionConfig := new(action.Configuration)
	settings := cli.New() // Requires helm-cli to be installed first
	err = actionConfig.Init(settings.RESTClientGetter(), "otterize-system", os.Getenv("HELM_DRIVER"), logrus.Debugf)
	s.Require().NoError(err)
	s.helmActionConfig = actionConfig
}

func (s *PostgresTestSuite) TestOtterizeKubernetesHelmInstall() {
	// Load Chart.yaml
	chart, err := loader.Load(OtterizeKubernetesChartPath)
	s.Require().NoError(err)

	installAction := action.NewInstall(s.helmActionConfig)
	installAction.Namespace = "otterize-system"
	installAction.ReleaseName = "otterize"
	installAction.CreateNamespace = true

	// Run helm install command
	results, err := installAction.Run(chart, chart.Values)
	s.Require().NoError(err)

	// Deploy client & database pods

}

func (s *PostgresTestSuite) TearDownTest() {

}

func (s *PostgresTestSuite) DeployPostgreSQLDatabasePod() {
	_, err := s.client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
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
						Name:                     "database",
						Image:                    "",
						Command:                  nil,
						Args:                     nil,
						WorkingDir:               "",
						Ports:                    nil,
						EnvFrom:                  nil,
						Env:                      nil,
						Resources:                corev1.ResourceRequirements{},
						ResizePolicy:             nil,
						RestartPolicy:            nil,
						VolumeMounts:             nil,
						VolumeDevices:            nil,
						LivenessProbe:            nil,
						ReadinessProbe:           nil,
						StartupProbe:             nil,
						Lifecycle:                nil,
						TerminationMessagePath:   "",
						TerminationMessagePolicy: "",
						ImagePullPolicy:          "",
						SecurityContext:          nil,
						Stdin:                    false,
						StdinOnce:                false,
						TTY:                      false,
					}},
					HostAliases:               nil,
					PriorityClassName:         "",
					Priority:                  nil,
					DNSConfig:                 nil,
					ReadinessGates:            nil,
					RuntimeClassName:          nil,
					EnableServiceLinks:        nil,
					PreemptionPolicy:          nil,
					Overhead:                  nil,
					TopologySpreadConstraints: nil,
					SetHostnameAsFQDN:         nil,
					OS:                        nil,
					HostUsers:                 nil,
					SchedulingGates:           nil,
					ResourceClaims:            nil,
				},
			},
		},
		Status: v1.DeploymentStatus{},
	}
}

func (s *PostgresTestSuite) ApplyIntents() {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "flink.apache.org/v1beta1",
			"kind":       "FlinkSessionJob",
			"metadata": map[string]interface{}{
				"name": "name of your application",
			},
			"spec": map[string]interface{}{
				"deploymentName": "Deployment name to which this FlinkSessionJob belongs",
				"job": map[string]interface{}{
					"jarURI":      "Your Jar file",
					"parallelism": 4,
					"upgradeMode": "upgrade mode",
					"state":       "running",
				},
			},
		},
	}
}

func TestPostgresEnforcementTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
