
# Parameters

## Deployment parameters
| Key                                   | Description                                 | Default |
|---------------------------------------|---------------------------------------------|---------|
| `deployment.spire`                    | Whether or not to deploy spire.             | `true`  |
| `deployment.spireIntegrationOperator` | Whether or not to deploy spire-integration. | `true`  |
| `deployment.intentsOperator`          | Whether or not to deploy intents-operator.  | `true`  |

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

## Intents operator parameters
All configurable parameters of intents-operator can be configured under the alias `intentsOperator`.
Further information about intents-operator parameters can be found [in the Intents Operator's helm chart](https://github.com/otterize/helm-charts/tree/main/intents-operator).

| Key                                                  | Description                                                    | Default          |
|------------------------------------------------------|----------------------------------------------------------------|------------------|
| `intentsOperator.autoGenerateTLSUsingSpireIntegration` | Use spire-integration to create TLS cert for intents-operator. | `true`           |
| `intentsOperator.operator.autoCreateNetworkPoliciesForExternalTraffic` | Automatically allow external traffic, if a new ClientIntents resource would result in blocking external (internet) traffic and there is an Ingress/Service resource indicating external traffic is expected. | `true`           |

## SPIRE parameters
All configurable parameters of SPIRE can be configured under the alias `spire`.
Further information about `SPIRE` parameters can be found [in SPIRE's helm chart](https://github.com/otterize/helm-charts/tree/main/spire).

## SPIRE integration operator parameters
All configurable parameters of the SPIRE integration operator can be configured under the alias `spireIntegrationOperator`.
Further information about SPIRE integration operator parameters can be found [in the SPIRE integration operator's chart](https://github.com/otterize/helm-charts/tree/main/spire-integration-operator).

## Resource configuration
| Component                  | Key                                  | Default  |
|----------------------------|--------------------------------------|----------|
| Intents Operator           | `intentsOperator.operator.resources` | `(none)` |
| Intents Operator - Watcher | `intentsOperator.watcher.resources`  | `(none)` |
| Spire Server               | `spire.server.resources`             | `(none)` |
| Spire Agent                | `spire.agent.resources`              | `(none)` |
| Spire Integration Operator | `spireIntegrationOperator.resources` | `(none)` |
