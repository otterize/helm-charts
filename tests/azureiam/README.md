

## Setup instructions

```shell
export SUBSCRIPTION_ID=ef54c90c-5351-4c8f-a126-16a6d789104f
export LOCATION="eastus"
export RESOURCE_GROUP="otterizeGitHubActionsResourceGroup"
export AKS_CLUSTER_NAME="otterizeAzureIAME2EAKSCluster"
export STORAGE_ACCOUNT_NAME=ottrazureiame2e

# Setup resource group
az group create --name $RESOURCE_GROUP --location $LOCATION

# Setup service principal for GitHub Actions
# Save the output json to store as a GitHub actions secret 
az ad sp create-for-rbac --name "otterizeGitHubActions" --role contributor \
  --scopes /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP --json-auth

# setup AKS cluster
az aks create -g $RESOURCE_GROUP -n $AKS_CLUSTER_NAME --node-count 1 --enable-oidc-issuer --enable-workload-identity --generate-ssh-keys
az aks get-credentials -n $AKS_CLUSTER_NAME -g $RESOURCE_GROUP

# Setup storage account
az storage account create \
 --name $STORAGE_ACCOUNT_NAME \
 --resource-group $RESOURCE_GROUP \
 --location $LOCATION

# Add the "storage blob data contributor" role assignment to the storage account
export GITHUB_APP_SP_ID=$(az ad sp list --display-name otterizeGitHubActions --query '[0].appId' -o tsv)
az role assignment create --role "Storage Blob Data Contributor" \
	--scope /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Storage/storageAccounts/$STORAGE_ACCOUNT_NAME \
	--assignee $GITHUB_APP_SP_ID

# [Optional] Add the "storage blob data contributor" role assignment to the storage account for the signed in user
export SIGNED_IN_USER_NAME=$(az ad signed-in-user show --query 'userPrincipalName' -o tsv)
az role assignment create --role "Storage Blob Data Contributor" \
	--scope /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Storage/storageAccounts/$STORAGE_ACCOUNT_NAME \
	--assignee $SIGNED_IN_USER_NAME


# apply the Otterize Azure IAM terraform module to setup Azure IAM identities for the Otterize operator
cd files/terraform
terraform init
terraform apply
```