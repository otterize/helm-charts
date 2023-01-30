
# Parameters

## Deployment parameters
| Key                              | Description                                    | Default |
|----------------------------------|------------------------------------------------|---------|
| `deployment.spire`               | Whether or not to deploy spire.                | `true`  |
| `deployment.credentialsOperator` | Whether or not to deploy credentials-operator. | `true`  |
| `deployment.intentsOperator`     | Whether or not to deploy intents-operator.     | `true`  |
| `deployment.networkMapper`       | Whether or not to deploy network-mapper.       | `true`  |

## Global parameters
These parameters are used by multiple charts, and must be kept the same for the correct functioning of the separate components.

| Key                                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Default         |
|----------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------|
| `global.spiffe.CASubject`              | The Subject that CA certificates should use (see below).	                                                                                                                                                                                                                                                                                                                                                                                                                                                               |                 |
| `global.spiffe.CASubject.country`      | Spire's CA certificates `Country` value.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | `"US"`          |
| `global.spiffe.CASubject.organization` | Spire's CA certificates `Organization` Value.                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | `"SPIRE"`       |
| `global.spiffe.trustDomain`            | The trust domain that spire will use.	                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | `"example.org"` |
| `global.spire.serverServiceName`       | Name of the kubernetes service that will be created for spire-server.                                                                                                                                                                                                                                                                                                                                                                                                                                                   |                 |
| `global.allowGetAllResources`          | If defined overrides `allowGetAllResources` in subcharts. Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. |                 |

## Cloud parameters
| Key                                             | Description                                     | Default  |
|-------------------------------------------------|-------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`     | Client ID for connecting to Otterize Cloud.     | `(none)` |
| `global.otterizeCloud.credentials.clientSecret` | Client secret for connecting to Otterize Cloud. | `(none)` |
| `global.otterizeCloud.apiAddress`               | Overrides Otterize Cloud default API address.   | `(none)` |

## Intents operator parameters
All configurable parameters of intents-operator can be configured under the alias `intentsOperator`.
Further information about intents-operator parameters can be found [in the Intents Operator's helm chart](https://github.com/otterize/helm-charts/tree/main/intents-operator).

| Key                                                                    | Description                                                                                                                                                                                                  | Default |
|------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `intentsOperator.autoGenerateTLSUsingCredentialsOperator`              | Use credentials-operator to create TLS cert for intents-operator.                                                                                                                                            | `true`  |
| `intentsOperator.operator.enableEnforcement`                           | If set to false, enforcement is disabled globally (both for network policies and Kafka ACL). If true, you may use the other flags for more granular enforcement settings                                     | `true`  |
| `intentsOperator.operator.enableNetworkPolicyCreation`                 | Whether the operator should create network policies according to the ClientIntents                                                                                                                           | `true`  |
| `intentsOperator.operator.enableKafkaACLCreation`                      | Whether the operator should create Kafka ACL rules according to the ClientIntents of type Kafka                                                                                                              | `true`  |
| `intentsOperator.operator.autoCreateNetworkPoliciesForExternalTraffic` | Automatically allow external traffic, if a new ClientIntents resource would result in blocking external (internet) traffic and there is an Ingress/Service resource indicating external traffic is expected. | `true`  |

## SPIRE parameters
All configurable parameters of SPIRE can be configured under the alias `spire`.
Further information about `SPIRE` parameters can be found [in SPIRE's helm chart](https://github.com/otterize/helm-charts/tree/main/spire).

## Credentials operator parameters
All configurable parameters of the Credentials operator can be configured under the alias `credentialsOperator`.
Further information about Credentials operator parameters can be found [in the Credentials operator's chart](https://github.com/otterize/helm-charts/tree/main/credentials-operator).

| Key                                    | Description                                                    | Default |
|----------------------------------------|----------------------------------------------------------------|---------|
| `credentialsOperator.useOtterizeCloud` | Use otterize cloud for certificate management instead of spire | `false` |

## Network mapper parameters
All configurable parameters of the network mapper can be configured under the alias `networkMapper`.
Further information about network mapper parameters can be found [in the network mapper's chart](https://github.com/otterize/helm-charts/tree/main/network-mapper).

## Resource configuration
| Component                  | Key                                  | Default  |
|----------------------------|--------------------------------------|----------|
| Intents Operator           | `intentsOperator.operator.resources` | `(none)` |
| Intents Operator - Watcher | `intentsOperator.watcher.resources`  | `(none)` |
| Spire Server               | `spire.server.resources`             | `(none)` |
| Spire Agent                | `spire.agent.resources`              | `(none)` |
| Credentials Operator       | `credentialsOperator.resources`      | `(none)` |
