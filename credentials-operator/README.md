# Parameters

## Global parameters
| Key                                  | Description                                                                                                                                 | Default |
|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.spire.serverServiceName`     | If deployed with SPIRE, this key specifies SPIRE-server's service name. You should use either this **OR** `spire.serverAddress` (not both). |         |
| `global.allowGetAllResources`        | If defined overrides `allowGetAllResources`.                                                                                                |         |
| `global.useOtterizeCloudCredentials` | Use Otterize Cloud for certificate management instead of SPIRE                                                                              | `false` |


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
| Key                                             | Description                                     | Default  |
|-------------------------------------------------|-------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`     | Client ID for connecting to Otterize Cloud.     | `(none)` |
| `global.otterizeCloud.credentials.clientSecret` | Client secret for connecting to Otterize Cloud. | `(none)` |
| `global.otterizeCloud.apiAddress`               | Overrides Otterize Cloud default API address.   | `(none)` |

## Common parameters
| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`  |
| `resources`            | Resources of the container                                                                                                                                                                                                                                                                                                                                                                                                                                    | `{}`    |
