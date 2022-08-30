
## Values

| Key                                                  | Type   | Default             | Description                                                                                             |
|------------------------------------------------------|--------|---------------------|---------------------------------------------------------------------------------------------------------|
| global.spire.serverServiceName                       | string | `"spire-server"`    | spire-server service name. Should be declared globally if both spire and spire-integration are deployed |
| spire.clusterName                                    | string | `"example-cluster"` |                                                                                                         |
| spire.server.rootCATTL                               | string | `"26280h"`          | determine root_ca TTL                                                                                   |
| spire.trustDomain                                    | string | `"example.org"`     |                                                                                                         |
| deployment.spire                                     | bool   | `true`              |                                                                                                         |
| deployment.spireIntegrationOperator                  | bool   | `true`              |                                                                                                         |
| deployment.intentsOperator                           | bool   | `true`              |                                                                                                         |
| intentsOperator.autoGenerateTLSUsingSpireIntegration | bool   | `true`              |                                                                                                         |
