package helm_tests

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm_tests/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var OtterizeKubernetesChartPath string

const (
	OtterizeNamespace       = "otterize-system"
	OtterizeHelmReleaseName = "otterize"
)

type BaseSuite struct {
	suite.Suite
	Client           *kubernetes.Clientset
	DynamicClient    *dynamic.DynamicClient
	HelmActionConfig *action.Configuration

	IntentsClient dynamic.NamespaceableResourceInterface
}

func (s *BaseSuite) GetKubeconfigPath() string {
	envPath := os.Getenv("KUBECONFIG")
	if envPath != "" {
		return envPath
	}

	homeDir, err := os.UserHomeDir()
	s.Require().NoError(err)

	return path.Join(homeDir, viper.GetString(config.KubeConfigPath))
}

func (s *BaseSuite) SetupSuite() {
	logrus.SetLevel(logrus.DebugLevel)
	kubeconfigPath := s.GetKubeconfigPath()
	logrus.WithField("kubeconfig", kubeconfigPath).Info("Using kubeconfig")
	s.Require().FileExists(kubeconfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	s.Require().NoError(err)

	client, err := kubernetes.NewForConfig(kubeConfig)
	s.Require().NoError(err)
	s.Client = client

	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	s.Require().NoError(err)
	s.DynamicClient = dynamicClient

	actionConfig := new(action.Configuration)
	settings := cli.New() // Requires helm-cli to be installed first
	err = actionConfig.Init(settings.RESTClientGetter(), OtterizeNamespace, os.Getenv("HELM_DRIVER"), logrus.Debugf)
	s.Require().NoError(err)
	s.HelmActionConfig = actionConfig

	s.IntentsClient = s.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "k8s.otterize.com",
		Version:  "v1alpha3",
		Resource: "clientintents",
	})
}

func (s *BaseSuite) InstallOtterizeHelmChart(values map[string]any) {
	// Load Chart.yaml
	chart, err := loader.Load(OtterizeKubernetesChartPath)
	s.Require().NoError(err)

	logrus.WithField("chart", chart.Metadata.Name).Info("Loaded helm chart")

	installAction := action.NewInstall(s.HelmActionConfig)
	installAction.Namespace = OtterizeNamespace
	installAction.ReleaseName = OtterizeHelmReleaseName
	installAction.CreateNamespace = true
	installAction.Wait = true
	installAction.Timeout = 2 * time.Minute

	// Run helm install command
	logrus.WithField("values", values).WithField("namespace", OtterizeNamespace).Info("Installing otterize helm chart")
	_, err = installAction.Run(chart, values)
	s.Require().NoError(err)
	logrus.Info("Otterize helm chart installed")
}

func (s *BaseSuite) UninstallOtterizeHelmChart(ctx context.Context) {
	logrus.Info("Uninstalling otterize helm chart")
	uninstallAction := action.NewUninstall(s.HelmActionConfig)
	_, err := uninstallAction.Run(OtterizeHelmReleaseName)
	s.Require().NoError(err)

	s.DeleteNamespace(ctx, OtterizeNamespace)
}

func (s *BaseSuite) DeleteNamespace(ctx context.Context, namespaceName string) {
	logrus.WithField("namespace", namespaceName).Info("Deleting namespace")
	err := s.Client.CoreV1().Namespaces().Delete(ctx, namespaceName, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return
	}
	s.Require().NoError(err)

	s.WaitForNamespaceDeletion(ctx, namespaceName)
}

func (s *BaseSuite) WaitForNamespaceDeletion(ctx context.Context, namespaceName string) {
	selector := fields.OneTermEqualSelector(metav1.ObjectNameField, namespaceName)
	watchOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	logger := logrus.WithField("namespace", namespaceName)

	watcher, err := s.Client.CoreV1().Namespaces().Watch(ctx, watchOptions)
	s.Require().NoError(err)
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		logger.WithField("type", event.Type).Debug("Namespace changed")

		switch event.Type {
		case watch.Deleted:
			logger.Info("Namespace deleted")
			return
		case watch.Error:
			s.Require().Failf("Unexpected namespace event type", "Unexpected namespace event type: %v", event.Type)
		default:
			continue
		}

	}
}

func (s *BaseSuite) CreateNamespace(ctx context.Context, namespaceName string) {
	logrus.WithField("namespace", namespaceName).Info("Creating namespace")
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	_, err := s.Client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *BaseSuite) CreateServiceAccount(ctx context.Context, namespaceName string, name string) {
	logrus.WithField("namespace", namespaceName).WithField("name", name).Info("Creating service account")
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := s.Client.CoreV1().ServiceAccounts(namespaceName).Create(ctx, serviceAccount, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *BaseSuite) CreateDeployment(ctx context.Context, deployment *appsv1.Deployment) {
	logrus.WithField("namespace", deployment.Namespace).WithField("deployment", deployment.Name).Info("Creating deployment")
	_, err := s.Client.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *BaseSuite) CreateService(ctx context.Context, service *corev1.Service) {
	logrus.WithField("namespace", service.Namespace).WithField("service", service.Name).Info("Creating service")
	_, err := s.Client.CoreV1().Services(service.Namespace).Create(ctx, service, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *BaseSuite) WaitForDeploymentAvailability(ctx context.Context, namespace string, deploymentName string) {
	selector := fields.OneTermEqualSelector(metav1.ObjectNameField, deploymentName)

	watchOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	watcher, err := s.Client.AppsV1().Deployments(namespace).Watch(ctx, watchOptions)
	s.Require().NoError(err)
	defer watcher.Stop()

	isDeploymentReady := func(dep *appsv1.Deployment) bool {
		_, readyConditionFound := lo.Find(dep.Status.Conditions, func(c appsv1.DeploymentCondition) bool {
			return c.Type == appsv1.DeploymentAvailable && c.Status == corev1.ConditionTrue
		})
		return readyConditionFound
	}

	for event := range watcher.ResultChan() {
		item := event.Object.(*appsv1.Deployment)
		logrus.WithField("name", item.Name).WithField("type", event.Type).Debug("Deployment changed")

		switch event.Type {
		case watch.Added:
		case watch.Modified:
			if isDeploymentReady(item) {
				return
			}
		case watch.Bookmark:
		case watch.Error:
		case watch.Deleted:
			s.Require().Failf("Unexpected deployment event type", "Unexpected deployment event type: %v", event.Type)
		}
	}
}

type LogLineMatcher func(line string) bool

func (s *BaseSuite) ReadPodLogsUntilSubstring(ctx context.Context, pod *corev1.Pod, substring string) {
	req := s.Client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Follow: true})
	logStream, err := req.Stream(ctx)
	s.Require().NoError(err)

	logger := logrus.WithField("pod", pod.Name).WithField("namespace", pod.Namespace)
	logger.Debugf("Reading pod logs searching for substring %s", substring)

	defer logStream.Close()

	reader := bufio.NewScanner(logStream)
	var line string
	for {
		select {
		case <-ctx.Done():
			logger.Errorf("Failed to match log substring: %s", substring)
			return
		default:
			for reader.Scan() {
				line = reader.Text()
				if strings.Contains(line, substring) {
					logger.Infof("Matched log line: %s", line)
					return
				}
			}
		}
	}
}

func (s *BaseSuite) FindPodByLabel(ctx context.Context, namespace string, labelSelector string) *corev1.Pod {
	pods, err := s.Client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	s.Require().NoError(err)

	s.Require().Lenf(pods.Items, 1, "Expected to find a single pod with label %s, found %d", labelSelector, len(pods.Items))

	return &pods.Items[0]
}

func (s *BaseSuite) DeleteClientIntents(ctx context.Context, namespaceName string, name string) {
	err := s.IntentsClient.Namespace(namespaceName).Delete(context.Background(), name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return
	}
	s.Require().NoError(err)

	s.WaitForClientIntentsDeletion(ctx, namespaceName, name)
}

func (s *BaseSuite) WaitForClientIntentsDeletion(ctx context.Context, namespaceName string, name string) {
	selector := fields.OneTermEqualSelector(metav1.ObjectNameField, name)
	watchOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	logger := logrus.WithField("namespace", namespaceName).WithField("name", name)

	watcher, err := s.IntentsClient.Namespace(namespaceName).Watch(ctx, watchOptions)
	s.Require().NoError(err)
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		logger.WithField("type", event.Type).Debug("ClientIntents changed")

		switch event.Type {
		case watch.Deleted:
			logger.Info("Namespace deleted")
			return
		case watch.Error:
			s.Require().Failf("Unexpected namespace event type", "Unexpected namespace event type: %v", event.Type)
		default:
			continue
		}

	}
}

func (s *BaseSuite) GetUnstructuredObject(resource any, gkv schema.GroupVersionKind) *unstructured.Unstructured {
	body, err := json.Marshal(resource)
	s.Require().NoError(err)

	u := unstructured.Unstructured{}
	err = u.UnmarshalJSON(body)
	s.Require().NoError(err)

	u.SetGroupVersionKind(gkv)
	return &u
}

func (s *BaseSuite) ApplyClientIntents(ctx context.Context, clientIntents v1alpha3.ClientIntents) {
	u := s.GetUnstructuredObject(clientIntents, clientIntents.GroupVersionKind())
	_, err := s.IntentsClient.Namespace(clientIntents.Namespace).Create(ctx, u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	OtterizeKubernetesChartPath = path.Join(path.Dir(filename), "..", "otterize-kubernetes")
	if _, err := os.Stat(OtterizeKubernetesChartPath); err != nil {
		panic(err)
	}
}
