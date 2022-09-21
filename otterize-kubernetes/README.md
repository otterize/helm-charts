
# Parameters

## Deployment parameters
| Key                                   | Description                                 | Default |
|---------------------------------------|---------------------------------------------|---------|
| `deployment.spire`                    | Whether or not to deploy spire.             | `true`  |
| `deployment.spireIntegrationOperator` | Whether or not to deploy spire-integration. | `true`  |
| `deployment.intentsOperator`          | Whether or not to deploy intents-operator.  | `true`  |

## Global parameters
| Key                                    | Description                                                           | Default         |
|----------------------------------------|-----------------------------------------------------------------------|-----------------|
| `global.spiffe.CASubject`              | The Subject that CA certificates should use (see below).	             |                 |
| `global.spiffe.CASubject.country`      | Spire's CA certificates `Country` value.                              | `"US"`          |
| `global.spiffe.CASubject.organization` | Spire's CA certificates `Organization` Value.                         | `"SPIRE"`       |
| `global.spiffe.trustDomain`            | The trust domain that spire will use.	                                | `"example.org"` |
| `global.spire.serverServiceName`       | Name of the kubernetes service that will be created for spire-server. |                 |
| `global.allowGetAllResources`          | If defined overrides `allowGetAllResources`.                          |                 |

## Intents-operator parameters
All configurable parameters of intents-operator can be configured under the alias IntentsOperator.
Further information about intents-operator parameters can be found [here](https://github.com/otterize/helm-charts/tree/main/intents-operator).

| Key                                                  | Description                                                    | Default          |
|------------------------------------------------------|----------------------------------------------------------------|------------------|
| intentsOperator.autoGenerateTLSUsingSpireIntegration | Use spire-integration to create TLS cert for intents-operator. | `true`           |

## Spire parameters
All configurable parameters of spire can be configured under the alias spire.
Further information about spire parameters can be found [here](https://github.com/otterize/helm-charts/tree/main/spire).

## Spire-integration-operator parameters
All configurable parameters of spire-integration-operator can be configured under the alias SpireIntegrationOperator.
Further information about spire-integration-operator parameters can be found [here](https://github.com/otterize/helm-charts/tree/main/spire-integration-operator).
