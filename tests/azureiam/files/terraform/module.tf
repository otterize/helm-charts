module "otterize-azure-iam" {
 source  = "otterize/otterize-azure-iam/azure"
 version = "1.0.1"
 azure_tenant_id = "f8b92b88-e477-41ad-a5af-079de8dc8210"
 azure_subscription_id = "ef54c90c-5351-4c8f-a126-16a6d789104f"
 azure_resource_group = "otterizeGitHubActionsResourceGroup"
 aks_cluster_name = "otterizeAzureIAME2EAKSCluster"
 otterize_deploy_namespace = "otterize-system"
}