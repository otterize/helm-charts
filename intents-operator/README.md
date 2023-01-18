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
| `operator.autoGenerateTLSUsingSpireIntegration`        | If true, adds the necessary pod annotations in order to integrate with spire-integration, and get tls certificate.                                                                                           | `true`             |
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

## Common parameters

| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`  |
