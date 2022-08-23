
## Values

| Key                                        | Type   | Default                           | Description           |
|--------------------------------------------|--------|-----------------------------------|-----------------------|
| affinity                                   | object | `{}`                              |                       |
| agent.image.pullPolicy                     | string | `"IfNotPresent"`                  |                       |
| agent.image.repository                     | string | `"gcr.io/spiffe-io/spire-agent"`  |                       |
| agent.image.tag                            | string | `""`                              |                       |
| agent.logLevel                             | string | `"INFO"`                          |                       |
| autoscaling.enabled                        | bool   | `false`                           |                       |
| autoscaling.maxReplicas                    | int    | `100`                             |                       |
| autoscaling.minReplicas                    | int    | `1`                               |                       |
| autoscaling.targetCPUUtilizationPercentage | int    | `80`                              |                       |
| fullnameOverride                           | string | `""`                              |                       |
| imagePullSecrets                           | list   | `[]`                              |                       |
| nameOverride                               | string | `""`                              |                       |
| nodeSelector                               | object | `{}`                              |                       |
| podAnnotations                             | object | `{}`                              |                       |
| podSecurityContext                         | object | `{}`                              |                       |
| replicaCount                               | int    | `1`                               |                       |
| resources                                  | object | `{}`                              |                       |
| securityContext                            | object | `{}`                              |                       |
| server.dataStorage.accessMode              | string | `"ReadWriteOnce"`                 |                       |
| server.dataStorage.enabled                 | bool   | `true`                            |                       |
| server.dataStorage.size                    | string | `"1Gi"`                           |                       |
| server.dataStorage.storageClass            | string | `nil`                             |                       |
| server.image.pullPolicy                    | string | `"IfNotPresent"`                  |                       |
| server.image.repository                    | string | `"gcr.io/spiffe-io/spire-server"` |                       |
| server.image.tag                           | string | `""`                              |                       |
| server.logLevel                            | string | `"INFO"`                          |                       |
| server.service.port                        | int    | `8081`                            |                       |
| server.service.type                        | string | `"ClusterIP"`                     |                       |
| serviceAccount.annotations                 | object | `{}`                              |                       |
| serviceAccount.create                      | bool   | `true`                            |                       |
| serviceAccount.name                        | string | `""`                              |                       |
| spire.clusterName                          | string | `"example-cluster"`               |                       |
| spire.server.root_ca_ttl                   | string | `"26280h"`                        | determine root_ca TTL |
| spire.trustDomain                          | string | `"example.org"`                   |                       |
| tolerations                                | list   | `[]`                              |                       |
