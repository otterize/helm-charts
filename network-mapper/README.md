# Parameters

## Mapper parameters
| Key                       | Description              | Default                                                               |
|---------------------------|--------------------------|-----------------------------------------------------------------------|
| `mapper.image.repository` | Mapper image repository  | `353146681200.dkr.ecr.us-east-1.amazonaws.com/otterize:mapper-latest` |
| `mapper.image.repository` | Mapper image pull policy | `Always`                                                              |


## Sniffer parameters
| Key                        | Description               | Default                                                                |
|----------------------------|---------------------------|------------------------------------------------------------------------|
| `sniffer.image.repository` | Sniffer image repository  | `353146681200.dkr.ecr.us-east-1.amazonaws.com/otterize:sniffer-latest` |
| `sniffer.image.repository` | Sniffer image pull policy | `Always`                                                               |


## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |

## Common parameters

| Key                    | Description                                    | Default |
|------------------------|------------------------------------------------|---------|
| `debug`                | Enable debug logs                              | `false` |
| `allowGetAllResources` | Gives get permission to watch on all resource. | `true`  |
