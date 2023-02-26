# Parameters

## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |

## Operator parameters
| Key                                                    | Description                                                                                                                                                                                                  | Default            |
|--------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| `operator.image.repository`                            | Intents Operator image repository.                                                                                                                                                                           | `otterize`         |
| `operator.image.image`                                 | Intents Operator image.                                                                                                                                                                                      | `intents-operator` |
| `operator.image.tag`                                   | Intents Operator image tag.                                                                                                                                                                                  | `latest`           |
| `operator.pullPolicy`                                  | Intents Operator image pull policy.                                                                                                                                                                          | `(none)`           |
| `operator.autoGenerateTLSUsingCredentialsOperator`        | If set to true, adds the necessary pod annotations in order to integrate with credentials-operator, and get tls certificate.                                                                                 | `false`            |
| `operator.enableEnforcement`                           | If set to false, enforcement is disabled globally (both for network policies and Kafka ACL). If true, you may use the other flags for more granular enforcement settings                                     | `true`             |
| `operator.enableNetworkPolicyCreation`                 | Whether the operator should create network policies according to ClientIntents                                                                                                                               | `true`             |
| `operator.enableKafkaACLCreation`                      | Whether the operator should create Kafka ACL rules according to ClientIntents of type Kafka                                                                                                                  | `true`             |
| `operator.autoCreateNetworkPoliciesForExternalTraffic` | Automatically allow external traffic, if a new ClientIntents resource would result in blocking external (internet) traffic and there is an Ingress/Service resource indicating external traffic is expected. | `true`             |
| `operator.resources`                                   | Resources override.                                                                                                                                                                                          |                    |

## Watcher parameters
| Key                        | Description                | Default                        |
|----------------------------|----------------------------|--------------------------------|
| `watcher.image.repository` | Watcher image repository.  | `otterize`                     |
| `watcher.image.image`      | Watcher image.             | `intents-operator-pod-watcher` |
| `watcher.image.tag`        | Watcher image tag.         | `latest`                       |
| `watcher.pullPolicy`       | Watcher image pull policy. | `(none)`                       |
| `watcher.resources`        | Watcher Resources.         |                                |

## Cloud parameters
| Key                                             | Description                                     | Default  |
|-------------------------------------------------|-------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`     | Client ID for connecting to Otterize Cloud.     | `(none)` |
| `global.otterizeCloud.credentials.clientSecret` | Client secret for connecting to Otterize Cloud. | `(none)` |
| `global.otterizeCloud.apiAddress`               | Overrides Otterize Cloud default API address.   | `(none)` |
| `global.otterizeCloud.apiExtraCAPEMSecret`      | The name of a secret containing extra root CA PEM file used to connect to Otterize Cloud. | `(none)` |

## Common parameters
| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`  |
