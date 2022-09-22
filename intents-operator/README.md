# Parameters

## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |


## Operator parameters

| Key                                             | Description                                                                                                        | Default            |
|-------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|--------------------|
| `operator.image.repository`                     | Intents Operator image repository.                                                                                 | `otterize`         |
| `operator.image.image`                          | Intents Operator image.                                                                                            | `intents-operator` |
| `operator.image.tag`                            | Intents Operator image tag.                                                                                        | `latest`           |
| `operator.pullPolicy`                           | Intents Operator image pull policy.                                                                                | `(none)`           |
| `operator.autoGenerateTLSUsingSpireIntegration` | If true, adds the necessary pod annotations in order to integrate with spire-integration, and get tls certificate. | `false`            |
| `operator.resources`                            | Resources override.                                                                                                |                    |


## Watcher parameters

| Key                        | Description                | Default                        |
|----------------------------|----------------------------|--------------------------------|
| `watcher.image.repository` | Watcher image repository.  | `otterize`                     |
| `watcher.image.image`      | Watcher image.             | `intents-operator-pod-watcher` |
| `watcher.image.tag`        | Watcher image tag.         | `latest`                       |
| `watcher.pullPolicy`       | Watcher image pull policy. | `(none)`                       |
| `watcher.resources`        | Watcher Resources.         |                                |

## Common parameters

| Key                    | Description                                    | Default |
|------------------------|------------------------------------------------|---------|
| `allowGetAllResources` | Gives get permission to watch on all resource. | `true`  |
