# Parameters

## Operator parameters

| Key                                             | Description                                                                                                        | Default                                                                         |
|-------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------|
| `operator.image.repository`                     | Operator image repository.                                                                                         | `353146681200.dkr.ecr.us-east-1.amazonaws.com/otterize:intents-operator-latest` |
| `operator.image.pullPolicy`                     | Operator image pull policy.                                                                                        | `Always`                                                                        |
| `operator.autoGenerateTLSUsingSpireIntegration` | If true, adds the necessary pod annotations in order to integrate with spire-integration, and get tls certificate. | `false`                                                                         |
| `operator.resources`                            | Operator resources.                                                                                                |                                                                                 |


## Watcher parameters

| Key                        | Description                | Default                                                                         |
|----------------------------|----------------------------|---------------------------------------------------------------------------------|
| `watcher.image.repository` | Watcher image repository.  | `353146681200.dkr.ecr.us-east-1.amazonaws.com/otterize:intents-operator-latest` |
| `watcher.image.pullPolicy` | Watcher image pull policy. | `Always`                                                                        |
| `watcher.resources`        | Watcher Resources.         |                                                                                 |

## Common parameters

| Key                    | Description                                    | Default |
|------------------------|------------------------------------------------|---------|
| `allowGetAllResources` | Gives get permission to watch on all resource. | `true`  |
