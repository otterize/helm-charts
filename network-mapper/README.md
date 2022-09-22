# Parameters

## Mapper parameters
| Key                       | Description              | Default                                                               |
|---------------------------|--------------------------|-----------------------------------------------------------------------|
| `mapper.image.repository` | Mapper image repository. | `otterize`                                                            |
| `mapper.image.image`      | Mapper image.            | `network-mapper`                                                      |
| `mapper.image.tag`        | Mapper image tag.        | `latest`                                                              |
| `mapper.pullPolicy`       | Mapper pull policy.      | `(none)`                                                              |


## Sniffer parameters
| Key                        | Description               | Default                  |
|----------------------------|---------------------------|--------------------------|
| `sniffer.image.repository` | Sniffer image repository. | `otterize`               |
| `sniffer.image.image`      | Sniffer image.            | `network-mapper-sniffer` |
| `sniffer.image.tag`        | Sniffer image tag.        | `latest`                 |
| `sniffer.pullPolicy`       | Sniffer pull policy.      | `(none)`                 |


## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |

## Common parameters

| Key                    | Description                                    | Default |
|------------------------|------------------------------------------------|---------|
| `debug`                | Enable debug logs                              | `false` |
| `allowGetAllResources` | Gives get permission to watch on all resource. | `true`  |
