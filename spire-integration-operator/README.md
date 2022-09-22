# Parameters

## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.spire.serverServiceName` | If deployed with spire, this key specifies spire-server's service name. You should use either this **OR** `spire.serverAddress` (not both). |         |
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |

## Spire parameters

| Key                   | Description                                                                                                    | Default                |
|-----------------------|----------------------------------------------------------------------------------------------------------------|------------------------|
| `spire.serverAddress` | Specify the spire-server's address. You should use either this OR `global.spire.serverServiceName` (not both). |                        |  
| `spire.socketsPath`   | Spire sockets path. The operator will expect to find agent.sock in the host-mounted folder                     | `"/run/spire/sockets"` |

## Operator parameters

| Key                         | Description                | Default                      |
|-----------------------------|----------------------------|------------------------------|
| `operator.image.repository` | Operator image repository. | `otterize`                   |
| `operator.image.image`      | Operator image.            | `spire-integration-operator` |
| `operator.image.tag`        | Operator image tag.        | `latest`                     |
| `operator.pullPolicy`       | Operator pull policy.      | `(none)`                     |

## Common parameters

| Key                    | Description                                   | Default |
|------------------------|-----------------------------------------------|---------|
| `allowGetAllResources` | Gives get permission to watch on all resource | `true`  |
| `resources`            | Resources of the container                    | `{}`    |

