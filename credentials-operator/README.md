# Parameters

## Global parameters

| Key                                                 | Description                                                                                                                                                                                                                                                                                                             | Default                              |
|-----------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| `global.certificateProvider`                        | What provider should be used to generate certificates/credentials - `"spire"`, `"otterize-cloud"` or `"cert-manager"`                                                                                                                                                                                                   | `"spire"`                            |
| `global.spire.serverServiceName`                    | If deployed with SPIRE, this key specifies SPIRE-server's service name. You should use either this **OR** `spire.serverAddress` (not both).                                                                                                                                                                             |                                      |
| `global.allowGetAllResources`                       | If defined overrides `allowGetAllResources`.                                                                                                                                                                                                                                                                            | `false`                              |
| `global.commonAnnotations`                          | Annotations to add to all deployed objects                                                                                                                                                                                                                                                                              | {}                                   |
| `global.commonLabels`                               | Labels to add to all deployed objects                                                                                                                                                                                                                                                                                   | {}                                   |
| `global.podAnnotations`                             | Annotations to add to all deployed pods                                                                                                                                                                                                                                                                                 | {}                                   |
| `global.podLabels`                                  | Labels to add to all deployed pods                                                                                                                                                                                                                                                                                      | {}                                   |
| `global.workloadNameOverrideAnnotationName`         | Which annotation to use (in the [service name resolution algorithm](https://docs.otterize.com/reference/service-identities#kubernetes-service-identity-resolution)) for setting a pod's service name, if not the default. Use this if you already have annotations on your pods that provide the correct workload name. | `intents.otterize.com/workload-name` |
| `global.aws.enabled`                                | Enable or disable AWS integration                                                                                                                                                                                                                                                                                       | `false`                              |
| `global.aws.eksClusterNameOverride`                 | EKS cluster name (overrides auto-detection)                                                                                                                                                                                                                                                                             | `(none)`                             |
| `global.azure.enabled`                              | Enable or disable Azure integration                                                                                                                                                                                                                                                                                     | `false`                              |
| `global.aws.useSoftDelete`                          | Use soft delete strategy (tag as deleted instead of actually delete) for AWS roles and policies                                                                                                                                                                                                                         | `false`                              |
| `global.gcp.enabled`                                | Enable or disable GCPs integration                                                                                                                                                                                                                                                                                      | `false`                              |
| `global.telemetry.enabled`                          | If set to `false`, all anonymous telemetries collection will be disabled                                                                                                                                                                                                                                                | `true`                               |
| `global.telemetry.usage.enabled`                    | If set to `false`, collection of anonymous telemetries on product usage will be disabled                                                                                                                                                                                                                                | `true`                               |
| `global.telemetry.errors.enabled`                   | If set to `false`, collection of anonymous telemetries on application crashes and errors will be disabled                                                                                                                                                                                                               | `true`                               |
| `global.telemetry.errors.endpointAddress`           | If set, overrides the default endpoint address for anonymous telemetries on application crashes and errors                                                                                                                                                                                                              | `(none)`                             |
| `global.telemetry.errors.stage`                     | If set, overrides the default stage for anonymous telemetries on application crashes and errors                                                                                                                                                                                                                         | `(none)`                             |
| `global.telemetry.errors.credentialsOperatorApiKey` | If set, overrides the default API key for anonymous telemetries on application crashes and errors                                                                                                                                                                                                                       | `(none)`                             |
| `global.openshift`                                  | Whether to configure and deploy SecurityContextConstraints that allow all components to run with minimal privileges on a default OpenShift installation.                                                                                                                                                                | `false`                              |

## SPIRE parameters

| Key                   | Description                                                                                                    | Default                |
|-----------------------|----------------------------------------------------------------------------------------------------------------|------------------------|
| `spire.serverAddress` | Specify the SPIRE-server's address. You should use either this OR `global.spire.serverServiceName` (not both). |                        |  
| `spire.socketsPath`   | SPIRE sockets path. The operator will expect to find agent.sock in the host-mounted folder                     | `"/run/spire/sockets"` |

## Operator parameters

| Key                                 | Description                          | Default                                                            |
|-------------------------------------|--------------------------------------|--------------------------------------------------------------------|
| `operator.repository`               | Operator image repository.           | `otterize`                                                         |
| `operator.image`                    | Operator image.                      | `credentials-operator`                                             |
| `operator.tag`                      | Operator image tag.                  | (pinned to latest version as of this Helm chart version's publish) |
| `operator.containerSecurityContext` | Security context for the containers. | `(consult values.yaml)`                                            |
| `operator.podSecurityContext`       | Security context for the pod.        | `(consult values.yaml)`                                            |
| `operator.pullPolicy`               | Operator pull policy.                | `(none)`                                                           |
| `operator.pullSecrets`              | Operator pull secrets.               | `(none)`                                                           |

## Cloud parameters

| Key                                                              | Description                                                                                                                                                                                  | Default  |
|------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`                      | Client ID for connecting to Otterize Cloud.                                                                                                                                                  | `(none)` |
| `global.otterizeCloud.credentials.clientSecret`                  | Client secret for connecting to Otterize Cloud.                                                                                                                                              | `(none)` |
| `global.otterizeCloud.credentials.clientSecretKeyRef.secretName` | If specified, the name of a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                                            | `(none)` |
| `global.otterizeCloud.credentials.clientSecretKeyRef.secretKey`  | If specified, the key for the clientSecret in a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                        | `(none)` |
| `global.otterizeCloud.apiAddress`                                | Overrides Otterize Cloud default API address.                                                                                                                                                | `(none)` |
| `global.otterizeCloud.apiExtraCAPEMSecret`                       | The name of a secret containing a single `CA.pem` file for an extra root CA used to connect to Otterize Cloud. The secret should be placed in the same namespace as the Otterize deployment. | `(none)` |

## cert-manager parameters

| Key                            | Description                                                                                                                       | Default |
|--------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|---------|
| `certManager.issuerName`       | The cert-manager Issuer (or ClusterIssuer if `useClusterIssuer` is set) to be used for certificate generation                     | `""`    |
| `certManager.useClusterIssuer` | Use ClusterIssuer. If false, looks for a namespace-scoped Issuer.                                                                 | `true`  |
| `certManager.autoApprove`      | Makes the credentials-operator auto-approve its CertificateRequests. Use when the cert-manager default auto-approver is disabled. | `false` |

## Common parameters

| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`  |
| `resources`            | Resources of the container                                                                                                                                                                                                                                                                                                                                                                                                                                    | `{}`    |
| `affinity`             | Pod affinity                                                                                                                                                                                                                                                                                                                                                                                                                                                  | `{}`    |
| `tolerations`          | Pod tolerations                                                                                                                                                                                                                                                                                                                                                                                                                                               | `[]`    |

## AWS integration parameters

| Key                                 | Description                                                                                     | Default  |
|-------------------------------------|-------------------------------------------------------------------------------------------------|----------|
| `aws.roleARN`                       | ARN of the AWS role the operator will use to access AWS.                                        | `(none)` |
| `global.aws.enabled`                | Enable or disable AWS integration                                                               | `false`  |
| `global.aws.eksClusterNameOverride` | EKS cluster name (overrides auto-detection)                                                     | `(none)` |
| `global.aws.useSoftDelete`          | Use soft delete strategy (tag as deleted instead of actually delete) for AWS roles and policies | `false`  |

## Azure integration parameters

| Key                                   | Description                                                            | Default  |
|---------------------------------------|------------------------------------------------------------------------|----------|
| `global.azure.userAssignedIdentityID` | ID of the user assigned identity used by the operator to access Azure. | `(none)` |
| `global.azure.subscriptionID`         | ID of the Azure subscription in which the AKS cluster is deployed.     | `(none)` |
| `global.azure.resoureceGroup`         | Name of the Azure resource group in which the AKS cluster is deployed. | `(none)` |
| `global.azure.aksClusterName`         | Name of the AKS cluster in which the operator is deployed.             | `(none)` |

## Credentials operator parameters

| Key                              | Description                                                                                                                   | Default |
|----------------------------------|-------------------------------------------------------------------------------------------------------------------------------|---------|
| `databaseSecretRotationInterval` | Interval in which secrets created by the credentials operator will be rotated. Valid time units are "ns", "ms", "s", "m", "h" | `8h`    |
| `enableSecretRotation`           | Whether periodic secret rotation is enabled                                                                                   | false   |
