# Parameters 

## Global parameters
| Key                                  | Description                                                                                                                                 | Default |
|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.spire.serverServiceName`     | If deployed with SPIRE, this key specifies SPIRE-server's service name. You should use either this **OR** `spire.serverAddress` (not both). |         |
| `global.allowGetAllResources`        | If defined overrides `allowGetAllResources`.                                                                                                |         |                                                                      | `false` |


## SPIRE parameters

| Key                   | Description                                                                                                    | Default                |
|-----------------------|----------------------------------------------------------------------------------------------------------------|------------------------|
| `spire.serverAddress` | Specify the SPIRE-server's address. You should use either this OR `global.spire.serverServiceName` (not both). |                        |  
| `spire.socketsPath`   | SPIRE sockets path. The operator will expect to find agent.sock in the host-mounted folder                     | `"/run/spire/sockets"` |

## Operator parameters

| Key                         | Description                | Default                      |
|-----------------------------|----------------------------|------------------------------|
| `operator.image.repository` | Operator image repository. | `otterize`                   |
| `operator.image.image`      | Operator image.            | `credentials-operator` |
| `operator.image.tag`        | Operator image tag.        | `latest`                     |
| `operator.pullPolicy`       | Operator pull policy.      | `(none)`                     |

## Cloud parameters
| Key                                                        | Description                                                                                                                                                                                  | Default  |
|------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `global.otterizeCloud.useCloudToGenerateTLSCredentials`    | Use Otterize Cloud for certificate management instead of SPIRE                                                                                                                               | `false`  |
| `global.otterizeCloud.credentials.clientId`                | Client ID for connecting to Otterize Cloud.                                                                                                                                                  | `(none)` |
| `global.otterizeCloud.credentials.clientSecret`            | Client secret for connecting to Otterize Cloud.                                                                                                                                              | `(none)` |
| `global.otterizeCloud.credentials.secretKeyRef.secretName` | If specified, the name of a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                                            | `(none)` |
| `global.otterizeCloud.credentials.secretKeyRef.secretKey`  | If specified, the key for the clientSecret in a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                        | `(none)` |
| `global.otterizeCloud.apiAddress`                          | Overrides Otterize Cloud default API address.                                                                                                                                                | `(none)` |
| `global.otterizeCloud.apiExtraCAPEMSecret`                 | The name of a secret containing a single `CA.pem` file for an extra root CA used to connect to Otterize Cloud. The secret should be placed in the same namespace as the Otterize deployment. | `(none)` |

## Common parameters
| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`  |
| `resources`            | Resources of the container                                                                                                                                                                                                                                                                                                                                                                                                                                    | `{}`    |
