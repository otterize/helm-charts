
## Values

| Key                                                  | Type   | Default             | Description           |
|------------------------------------------------------|--------|---------------------|-----------------------|
| spire.clusterName                                    | string | `"example-cluster"` |                       |
| spire.server.rootCATTL                               | string | `"26280h"`          | determine root_ca TTL |
| spire.trustDomain                                    | string | `"example.org"`     |                       |
| deployment.spire                                     | bool   | `true`              |                       |
| deployment.spireIntegrationOperator                  | bool   | `true`              |                       |
| deployment.intentsOperator                           | bool   | `true`              |                       |
| intentsOperator.autoGenerateTLSUsingSpireIntegration | bool   | `true`              |                       |
