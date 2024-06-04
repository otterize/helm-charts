package databases

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm_tests/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type BaseSuite struct {
	suite.Suite
	Client             *kubernetes.Clientset
	IntentsClient      dynamic.NamespaceableResourceInterface
	PGServerConfClient dynamic.NamespaceableResourceInterface
	HelmActionConfig   *action.Configuration
}

func (s *BaseSuite) SetupSuite() {
	homeDir, err := os.UserHomeDir()
	s.Require().NoError(err)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/%s", homeDir, viper.GetString(config.KubeConfigPath)))
	s.Require().NoError(err)

	//
	client, err := kubernetes.NewForConfig(kubeConfig)
	s.Require().NoError(err)
	s.Client = client

	actionConfig := new(action.Configuration)
	settings := cli.New() // Requires helm-cli to be installed first
	err = actionConfig.Init(settings.RESTClientGetter(), "otterize-system", os.Getenv("HELM_DRIVER"), logrus.Debugf)
	s.Require().NoError(err)
	s.HelmActionConfig = actionConfig

	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	s.Require().NoError(err)
	s.IntentsClient = dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "k8s.otterize.com",
		Version:  "v1alpha3",
		Resource: "clientintents",
	})

	s.PGServerConfClient = dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "k8s.otterize.com",
		Version:  "v1alpha3",
		Resource: "postgresqlserverconfigs",
	})
}
