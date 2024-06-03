package azureiam

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/otterize/intents-operator/src/shared/agentutils"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	tests "helm_tests"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
	"testing"
	"time"
)

type AzureConfig struct {
	SubscriptionID     string `env:"AZURE_SUBSCRIPTION_ID"`
	Location           string `env:"AZURE_LOCATION"`
	ResourceGroup      string `env:"AZURE_RESOURCE_GROUP"`
	AKSClusterName     string `env:"AZURE_AKS_CLUSTER_NAME"`
	StorageAccountName string `env:"AZURE_STORAGE_ACCOUNT_NAME"`

	OtterizeOperatorUserAssignedIdentityClientID string `env:"AZURE_OTTERIZE_OPERATOR_USER_ASSIGNED_IDENTITY_CLIENT_ID"`
}

const (
	OtterizeKubernetesChartPath = "../../otterize-kubernetes"
	OtterizeNamespace           = "otterize-system"
	azBlobFileName              = "hello.txt"
	clientAppNamespaceName      = "otterize-tutorial-azure-iam"
	clientAppServiceAccountName = "client"
	clientAppDeploymentName     = "client"

	// constants taken from intents-operator/src/shared/azureagent/identities.go
	maxManagedIdentityLength   = 128
	maxFederatedIdentityLength = 120
)

const clientContainerCommandArgs = `while true;
do
	echo;
	echo 'Client - The time is:' $(date);
	echo;
	if [[ -z "$AZURE_CLIENT_ID" ]]; then echo "Azure client ID not set";
	else
	  echo 'Logging in using federated identity credentials';
	  az login -o table --federated-token $(cat $AZURE_FEDERATED_TOKEN_FILE) --service-principal -u $AZURE_CLIENT_ID -t $AZURE_TENANT_ID;
	  echo;
	  echo 'Listing storage blob container' $AZURE_STORAGE_CONTAINER 'in storage account' $AZURE_STORAGE_ACCOUNT;
	  az storage blob list --container $AZURE_STORAGE_CONTAINER --account-name $AZURE_STORAGE_ACCOUNT --auth-mode login -o table;
	  echo;
	fi;
	sleep 5;
done`

// function taken from intents-operator/src/shared/azureagent/identities.go
func generateUserAssignedIdentityName(namespace string, accountName string, aksClusterName string) string {
	fullName := fmt.Sprintf("ottr-uai-%s-%s-%s", namespace, accountName, aksClusterName)
	return agentutils.TruncateHashName(fullName, maxManagedIdentityLength)
}

// function taken from intents-operator/src/shared/azureagent/identities.go
func generateFederatedIdentityCredentialsName(namespace string, accountName string, aksClusterName string) string {
	fullName := fmt.Sprintf("ottr-fic-%s-%s-%s", namespace, accountName, aksClusterName)
	return agentutils.TruncateHashName(fullName, maxFederatedIdentityLength)
}

type AzureIAMTestSuite struct {
	tests.BaseSuite
	conf AzureConfig

	// Azure clients
	credentials                        *azidentity.DefaultAzureCredential
	storageClientFactory               *armstorage.ClientFactory
	accountsClient                     *armstorage.AccountsClient
	blobContainersClient               *armstorage.BlobContainersClient
	userAssignedIdentitiesClient       *armmsi.UserAssignedIdentitiesClient
	federatedIdentityCredentialsClient *armmsi.FederatedIdentityCredentialsClient
	azBlobClient                       *azblob.Client
}

func (s *AzureIAMTestSuite) installOtterizeForAzureIAM() {
	// Load Chart.yaml
	chart, err := loader.Load(OtterizeKubernetesChartPath)
	s.Require().NoError(err)

	logrus.WithField("chart", chart.Metadata.Name).Info("Loaded helm chart")

	installAction := action.NewInstall(s.HelmActionConfig)
	installAction.Namespace = OtterizeNamespace
	installAction.ReleaseName = "otterize"
	installAction.CreateNamespace = true
	installAction.Wait = true
	installAction.Timeout = 2 * time.Minute

	// Run helm install command
	values := map[string]any{
		"global": map[string]any{
			"azure": map[string]any{
				"enabled":                true,
				"subscriptionID":         s.conf.SubscriptionID,
				"resourceGroup":          s.conf.ResourceGroup,
				"aksClusterName":         s.conf.AKSClusterName,
				"userAssignedIdentityID": s.conf.OtterizeOperatorUserAssignedIdentityClientID,
			},
			"deployment": map[string]any{
				"networkMapper": false,
			},
		},
	}

	logrus.WithField("values", values).WithField("namespace", OtterizeNamespace).Info("Installing otterize helm chart")
	_, err = installAction.Run(chart, values)
	s.Require().NoError(err)
	logrus.Info("Otterize helm chart installed")
}

func (s *AzureIAMTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	s.Require().NoError(godotenv.Load("azure-account.env"))
	s.Require().NoError(env.Parse(&s.conf))
	s.initAzureAgent()

	s.installOtterizeForAzureIAM()
}

func (s *AzureIAMTestSuite) uninstallOtterize() {
	logrus.Info("Uninstalling otterize helm chart")
	uninstallAction := action.NewUninstall(s.HelmActionConfig)
	_, err := uninstallAction.Run("otterize")
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) deleteOtterizeNamespace(ctx context.Context) {
	logrus.WithField("namespace", OtterizeNamespace).Info("Deleting otterize namespace")
	err := s.Client.CoreV1().Namespaces().Delete(ctx, OtterizeNamespace, metav1.DeleteOptions{})
	s.Require().NoError(err)

	s.waitForNamespaceDeletion(ctx, OtterizeNamespace)
}

func (s *AzureIAMTestSuite) TearDownSuite() {
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	//defer cancel()
	//s.uninstallOtterize()
	//s.deleteOtterizeNamespace(ctx)
}

func (s *AzureIAMTestSuite) cleanupClientApp(ctx context.Context) {
	logrus.WithField("namespace", clientAppNamespaceName).Info("Deleting client app namespace")
	err := s.Client.CoreV1().Namespaces().Delete(ctx, clientAppNamespaceName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		s.Require().NoError(err)
	}

	s.waitForNamespaceDeletion(ctx, clientAppNamespaceName)
}

func (s *AzureIAMTestSuite) waitForNamespaceDeletion(ctx context.Context, namespace string) {
	selector := fields.OneTermEqualSelector(metav1.ObjectNameField, namespace)
	watchOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	watcher, err := s.Client.CoreV1().Namespaces().Watch(ctx, watchOptions)
	s.Require().NoError(err)
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		item := event.Object.(*v1.Namespace)
		logrus.WithField("name", item.Name).WithField("type", event.Type).Info("Namespace changed")

		switch event.Type {
		case watch.Deleted:
			logrus.WithField("namespace", item.Name).Info("Namespace deleted")
			return
		case watch.Error:
			s.Require().Failf("Unexpected namespace event type", "Unexpected namespace event type: %v", event.Type)
		default:
			continue
		}

	}
}

func (s *AzureIAMTestSuite) TearDownTest() {
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
	//defer cancel()
	//
	//s.cleanupClientApp(ctx)
}

func (s *AzureIAMTestSuite) initAzureAgent() {
	var err error

	s.credentials, err = azidentity.NewDefaultAzureCredential(nil)
	s.Require().NoError(err)

	s.storageClientFactory, err = armstorage.NewClientFactory(s.conf.SubscriptionID, s.credentials, nil)
	s.Require().NoError(err)

	s.accountsClient = s.storageClientFactory.NewAccountsClient()
	s.blobContainersClient = s.storageClientFactory.NewBlobContainersClient()

	armmsiClientFactory, err := armmsi.NewClientFactory(s.conf.SubscriptionID, s.credentials, nil)
	s.Require().NoError(err)
	s.userAssignedIdentitiesClient = armmsiClientFactory.NewUserAssignedIdentitiesClient()
	s.federatedIdentityCredentialsClient = armmsiClientFactory.NewFederatedIdentityCredentialsClient()

	storageAccountURL := fmt.Sprintf("https://%s.blob.core.windows.net", s.conf.StorageAccountName)
	s.azBlobClient, err = azblob.NewClient(storageAccountURL, s.credentials, nil)
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) uploadTestBlobFile(ctx context.Context, containerName string) {
	logrus.WithField("container", containerName).Info("Creating Azure Blob Storage container")
	_, err := s.azBlobClient.CreateContainer(ctx, containerName, nil)
	s.Require().NoError(err)

	blobName := azBlobFileName
	data := []byte("Hello, Azure integration!")

	logrus.WithField("blob", blobName).Info("Uploading test blob file")
	_, err = s.azBlobClient.UploadBuffer(ctx, containerName, blobName, data, &azblob.UploadBufferOptions{})
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) createClientAppNamespace(ctx context.Context) {
	logrus.WithField("namespace", clientAppNamespaceName).Info("Creating client app namespace")
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: clientAppNamespaceName,
		},
	}

	_, err := s.Client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) createClientAppServiceAccount(ctx context.Context) {
	logrus.WithField("serviceAccount", clientAppServiceAccountName).Info("Creating client app service account")
	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clientAppServiceAccountName,
			Namespace: clientAppNamespaceName,
		},
	}

	_, err := s.Client.CoreV1().ServiceAccounts(serviceAccount.Namespace).Create(ctx, serviceAccount, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) createClientAppDeployment(ctx context.Context, storageAccountName string, storageContainerName string) {
	logrus.WithField("deployment", clientAppDeploymentName).Info("Creating client app deployment")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clientAppDeploymentName,
			Namespace: clientAppNamespaceName,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "client",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                         "client",
						"azure.workload.identity/use": "true",
						"credentials-operator.otterize.com/create-azure-workload-identity": "true",
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: clientAppServiceAccountName,
					Containers: []v1.Container{
						{
							Name:    "client",
							Image:   "mcr.microsoft.com/azure-cli",
							Command: []string{"/bin/sh", "-c", "--"},
							Env: []v1.EnvVar{
								{
									Name:  "AZURE_STORAGE_ACCOUNT",
									Value: storageAccountName,
								},
								{
									Name:  "AZURE_STORAGE_CONTAINER",
									Value: storageContainerName,
								},
							},
							Args: []string{clientContainerCommandArgs},
						},
					},
				},
			},
		},
	}

	_, err := s.Client.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	s.Require().NoError(err)
}

func (s *AzureIAMTestSuite) deployAzureBlobStorageClientApp(ctx context.Context, storageContainerName string) {
	logrus.WithField("namespace", clientAppNamespaceName).Info("Deploying Azure Blob Storage client app")
	s.createClientAppNamespace(ctx)
	s.createClientAppServiceAccount(ctx)
	s.createClientAppDeployment(ctx, s.conf.StorageAccountName, storageContainerName)

	s.waitForDeploymentAvailability(ctx, clientAppNamespaceName, clientAppDeploymentName)
	logrus.WithField("namespace", clientAppNamespaceName).Info("Client app deployment is ready")
}

func (s *AzureIAMTestSuite) waitForDeploymentAvailability(ctx context.Context, namespace string, deploymentName string) {
	selector := fields.OneTermEqualSelector(metav1.ObjectNameField, deploymentName)

	watchOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	watcher, err := s.Client.AppsV1().Deployments(namespace).Watch(ctx, watchOptions)
	s.Require().NoError(err)
	defer watcher.Stop()

	isDeploymentReady := func(dep *appsv1.Deployment) bool {
		_, readyConditionFound := lo.Find(dep.Status.Conditions, func(c appsv1.DeploymentCondition) bool {
			return c.Type == appsv1.DeploymentAvailable && c.Status == v1.ConditionTrue
		})
		return readyConditionFound
	}

	for event := range watcher.ResultChan() {
		item := event.Object.(*appsv1.Deployment)
		logrus.WithField("name", item.Name).WithField("type", event.Type).Info("Deployment changed")

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

func (s *AzureIAMTestSuite) readPodLogsUntilLine(ctx context.Context, pod *v1.Pod, matchLine LogLineMatcher) {
	req := s.Client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{Follow: true})
	logStream, err := req.Stream(ctx)
	s.Require().NoError(err)

	defer logStream.Close()

	logger := logrus.WithField("pod", pod.Name).WithField("namespace", pod.Namespace)
	reader := bufio.NewScanner(logStream)
	var line string
	for {
		select {
		case <-ctx.Done():
			break
		default:
			for reader.Scan() {
				line = reader.Text()
				logger.Debug(line)
				if matchLine(line) {
					return
				}
			}
		}
	}
}

func (s *AzureIAMTestSuite) findClientAppPod(ctx context.Context) *v1.Pod {
	pods, err := s.Client.CoreV1().Pods(clientAppNamespaceName).List(ctx, metav1.ListOptions{LabelSelector: "app=client"})
	s.Require().NoError(err)

	s.Require().Lenf(pods.Items, 1, "Expected to find a single pod with label app=client, found %d", len(pods.Items))

	pod := pods.Items[0]
	return &pod
}

func (s *AzureIAMTestSuite) waitUntilClientAppLogInUsingFederatedIdentityCredentials(ctx context.Context) {
	pod := s.findClientAppPod(ctx)
	logrus.WithField("pod", pod.Name).Info("Waiting for client app to log in using federated identity credentials")
	s.readPodLogsUntilLine(ctx, pod, func(line string) bool {
		return strings.Contains(line, "Logging in using federated identity credentials")
	})
}

func (s *AzureIAMTestSuite) waitUntilClientAppLogsListingStorageContainer(ctx context.Context, storageContainerName string) {
	pod := s.findClientAppPod(ctx)
	logrus.WithField("pod", pod.Name).Info("Waiting for client app to list storage container")
	expectedLine := fmt.Sprintf("Listing storage blob container %s in storage account %s", storageContainerName, s.conf.StorageAccountName)
	s.readPodLogsUntilLine(ctx, pod, func(line string) bool {
		return strings.Contains(line, expectedLine)
	})
}

func (s *AzureIAMTestSuite) waitUntilClientAppAllowedBlobAccess(ctx context.Context) {
	pod := s.findClientAppPod(ctx)
	logrus.WithField("pod", pod.Name).Info("Waiting for client app to successfully list blob container content")
	s.readPodLogsUntilLine(ctx, pod, func(line string) bool {
		return strings.Contains(line, azBlobFileName)
	})
}

func (s *AzureIAMTestSuite) ensureAzureWorkloadIdentityCreated(ctx context.Context) (uai armmsi.Identity, fic armmsi.FederatedIdentityCredential) {
	uaiName := generateUserAssignedIdentityName(clientAppNamespaceName, clientAppServiceAccountName, s.conf.AKSClusterName)
	uaiResponse, err := s.userAssignedIdentitiesClient.Get(ctx, s.conf.ResourceGroup, uaiName, nil)
	s.Require().NoError(err)
	uai = uaiResponse.Identity

	logrus.WithField("userAssignedIdentity", uai.Name).Info("User assigned identity found")

	ficName := generateFederatedIdentityCredentialsName(clientAppNamespaceName, clientAppServiceAccountName, s.conf.AKSClusterName)
	ficResponse, err := s.federatedIdentityCredentialsClient.Get(ctx, s.conf.ResourceGroup, uaiName, ficName, nil)
	s.Require().NoError(err)

	fic = ficResponse.FederatedIdentityCredential

	logrus.WithField("federatedIdentityCredentials", fic.Name).Info("Federated identity credentials found")

	return uai, fic
}

func (s *AzureIAMTestSuite) ensureServiceAccountLabeledWithAzureWorkloadIdentityClientID(ctx context.Context, uai armmsi.Identity) {
	serviceAccount, err := s.Client.CoreV1().ServiceAccounts(clientAppNamespaceName).Get(ctx, clientAppServiceAccountName, metav1.GetOptions{})
	s.Require().NoError(err)

	value, ok := serviceAccount.Annotations["azure.workload.identity/client-id"]
	s.Require().True(ok, "Expected to find annotation azure.workload.identity/client-id on service account")
	s.Require().Equal(*uai.Properties.ClientID, value, "Expected service account annotation azure.workload.identity/client-id to match user assigned identity client ID")

	logrus.WithField("serviceAccount", serviceAccount.Name).Info("Service account annotated with Azure workload identity client ID")
}

func (s *AzureIAMTestSuite) applyClientIntents(ctx context.Context, storageContainerName string) {
	logrus.Info("Applying client intents")
	storageContainerScope := fmt.Sprintf("/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s", s.conf.StorageAccountName, storageContainerName)

	u := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "k8s.otterize.com/v1alpha3",
			"kind":       "ClientIntents",
			"metadata": map[string]any{
				"name":      "client",
				"namespace": clientAppNamespaceName,
			},
			"spec": map[string]any{
				"service": map[string]any{
					"name": clientAppServiceAccountName,
				},
				"calls": []map[string]any{
					{
						"type": "azure",
						"name": storageContainerScope,
						"azureRoles": []string{
							"Storage Blob Data Contributor",
						},
					},
				},
			},
		},
	}

	gvk := u.GroupVersionKind()

	resource := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "clientintents",
	}

	_, err := s.DynamicClient.Resource(resource).Namespace(clientAppNamespaceName).Create(ctx, u, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// TestOtterizeKubernetesForAzureDemoFlow tests the end-to-end flow of deploying an Azure Blob Storage client app in an AKS cluster managed by Otterize with Azure integration.
// This test follows the tutorial flow described here: https://docs.otterize.com/features/azure-iam/tutorials/azure-iam-aks
func (s *AzureIAMTestSuite) TestOtterizeKubernetesForAzureDemoFlow() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
	defer cancel()
	//  Create an Azure Blob Storage account & container
	containerName := fmt.Sprintf("test%d", time.Now().Unix())
	s.uploadTestBlobFile(ctx, containerName)

	// Deploy the sample client
	s.deployAzureBlobStorageClientApp(ctx, containerName)

	// An Azure workload identity was created for the client pod
	uai, _ := s.ensureAzureWorkloadIdentityCreated(ctx)

	// The Kubernetes ServiceAccount was annotated with the workload identity ID
	s.ensureServiceAccountLabeledWithAzureWorkloadIdentityClientID(ctx, uai)

	// View logs for the client - Azure client ID is set, but no subscriptions found
	s.waitUntilClientAppLogInUsingFederatedIdentityCredentials(ctx)

	// Apply intents to create the necessary IAM role assignments
	s.applyClientIntents(ctx, containerName)

	// The client can now list files in the Azure Blob Storage container!
	s.waitUntilClientAppLogsListingStorageContainer(ctx, containerName)
	s.waitUntilClientAppAllowedBlobAccess(ctx)
}

func TestAzureIAMTestSuite(t *testing.T) {
	suite.Run(t, new(AzureIAMTestSuite))
}
