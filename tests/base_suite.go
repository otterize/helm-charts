package tests

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm_tests/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
)

type BaseSuite struct {
	suite.Suite
	Client           *kubernetes.Clientset
	DynamicClient    *dynamic.DynamicClient
	HelmActionConfig *action.Configuration
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
	err = actionConfig.Init(settings.RESTClientGetter(), "otterize-system", os.Getenv("HELM_DRIVER"), logrus.Debugf)
	s.Require().NoError(err)
	s.HelmActionConfig = actionConfig
}
