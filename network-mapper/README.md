# Parameters

## Mapper parameters

| Key                               | Description                                                                                     | Default                                                            |
|-----------------------------------|-------------------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| `mapper.repository`               | Mapper image repository.                                                                        | `otterize`                                                         |
| `mapper.image`                    | Mapper image.                                                                                   | `network-mapper`                                                   |
| `mapper.tag`                      | Mapper image tag.                                                                               | (pinned to latest version as of this Helm chart version's publish) |
| `mapper.containerSecurityContext` | Security context for the containers.                                                            | `(consult values.yaml)`                                            |
| `mapper.podSecurityContext`       | Security context for the pod.                                                                   | `(consult values.yaml)`                                            |
| `mapper.pullPolicy`               | Mapper pull policy.                                                                             | `(none)`                                                           |
| `mapper.pullSecrets`              | Mapper pull secrets.                                                                            | `(none)`                                                           |
| `mapper.resources`                | Resources override.                                                                             | `(none)`                                                           |
| `mapper.affinity`                 | Pod affinity.                                                                                   | `{}`                                                               |
| `mapper.nodeSelector`             | Node selector for the mapper.                                                                   | `{}`                                                               |
| `mapper.tolerations`              | Pod tolerations.                                                                                | `[]`                                                               |
| `mapper.uploadIntervalSeconds`    | Interval for uploading data to cloud                                                            | `60`                                                               |
| `mapper.excludeNamespaces`        | Namespaces excluded from reporting                                                              | `[istio-system]`                                                   |
| `mapper.extraEnvVars`             | List of extra env vars for the mapper, formatted as in the Kubernetes PodSpec (name and value). | `(none)`                                                           |

## Internet-facing traffic reporting

| Key                                    | Description                                                                                                                 | Default |
|----------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|---------|
| `enableInternetFacingTrafficReporting` | Whether to report internet-facing traffic to Otterize Cloud. This is a temporary flag that will soon be enabled by default. | `false` |

## OpenTelemetry exporter parameters

| Key                        | Description                                                                                                                                                                                                     | Default                              |
|----------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| `opentelemetry.enable`     | Whether to enable the OpenTelemetry exporter, which exports Grafana Tempo-style metrics for your network map. Configure the OpenTelemetry SDK using `mapper.extraEnvVars` (e.g. `OTEL_EXPORTER_OTLP_ENDPOINT`). | `false`                              |
| `opentelemetry.metricName` | The name of the OpenTelemetry metric name exported for the Grafana Tempo-style metric.                                                                                                                          | `traces_service_graph_request_total` |

## Sniffer parameters

| Key                                | Description                          | Default                                                            |
|------------------------------------|--------------------------------------|--------------------------------------------------------------------|
| `sniffer.enable`                   | Enable sniffer deployment.           | `true`                                                             |
| `sniffer.repository`               | Sniffer image repository.            | `otterize`                                                         |
| `sniffer.image`                    | Sniffer image.                       | `network-mapper-sniffer`                                           |
| `sniffer.containerSecurityContext` | Security context for the containers. | `(consult values.yaml)`                                            |
| `sniffer.podSecurityContext`       | Security context for the pods.       | `(consult values.yaml)`                                            |
| `sniffer.tag`                      | Sniffer image tag.                   | (pinned to latest version as of this Helm chart version's publish) |
| `sniffer.pullPolicy`               | Sniffer pull policy.                 | `(none)`                                                           |
| `sniffer.pullSecrets`              | Sniffer pull secrets.                | `(none)`                                                           |
| `sniffer.resources`                | Resources override.                  | `(none)`                                                           |
| `sniffer.affinity`                 | Sniffer's pod affinity.              | `{}`                                                               |
| `sniffer.nodeSelector`             | Sniffer's pod node selector.         | `{}`                                                               |
| `sniffer.tolerations`              | Sniffer's pod tolerations.           | `[]`                                                               |
| `sniffer.priorityClassName`        | Set priorityClassName.               | `(none)`                                                           |

## Kafka watcher parameters

| Key                                     | Description                                                 | Default                                                            |
|-----------------------------------------|-------------------------------------------------------------|--------------------------------------------------------------------|
| `kafkawatcher.enable`                   | Enable Kafka watcher deployment (beta).                     | `false`                                                            |
| `kafkawatcher.repository`               | Kafka watcher image repository.                             | `otterize`                                                         |
| `kafkawatcher.image`                    | Kafka watcher image.                                        | `network-mapper-kafka-watcher`                                     |
| `kafkawatcher.containerSecurityContext` | Security context for the containers.                        | `(consult values.yaml)`                                            |
| `kafkawatcher.podSecurityContext`       | Security context for the pod.                               | `(consult values.yaml)`                                            |
| `kafkawatcher.tag`                      | Kafka watcher image tag.                                    | (pinned to latest version as of this Helm chart version's publish) |
| `kafkawatcher.pullPolicy`               | Kafka watcher pull policy.                                  | `(none)`                                                           |
| `kafkawatcher.pullSecrets`              | Kafka watcher pull secrets.                                 | `(none)`                                                           |
| `kafkawatcher.resources`                | Resources override.                                         | `(none)`                                                           |
| `kafkawatcher.affinity`                 | Pod affinity.                                               | `{}`                                                               |
| `kafkawatcher.nodeSelector`             | Node selector for the Kafka watcher.                        | `{}`                                                               |
| `kafkawatcher.tolerations`              | Pod tolerations.                                            | `[]`                                                               |
| `kafkawatcher.kafkaServers`             | Kafka servers to watch, specified as `pod.namespace` items. | `(none)`                                                           |

## DNS visibility parameters

Deployed only when `aws.visibility.enabled` is set to `true`.

| Key                                      | Description                          | Default                 |
|------------------------------------------|--------------------------------------|-------------------------|
| `visibilitydns.repository`               | Image repository.                    | `coredns`               |
| `visibilitydns.image`                    | Image.                               | `coredns`               |
| `visibilitydns.containerSecurityContext` | Security context for the containers. | `(consult values.yaml)` |
| `visibilitydns.podSecurityContext`       | Security context for the pod.        | `(consult values.yaml)` |
| `visibilitydns.tag`                      | Image tag.                           | `latest`                |
| `visibilitydns.pullPolicy`               | Pull policy.                         | `(none)`                |
| `visibilitydns.pullSecrets`              | Pull secrets.                        | `(none)`                |
| `visibilitydns.resources`                | Resources override.                  | `(none)`                |

## Cloud parameters

| Key                                                              | Description                                                                                                                                                                                  | Default  |
|------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`                      | Client ID for connecting to Otterize Cloud.                                                                                                                                                  | `(none)` |
| `global.otterizeCloud.credentials.clientSecret`                  | Client secret for connecting to Otterize Cloud.                                                                                                                                              | `(none)` |
| `global.otterizeCloud.credentials.clientSecretKeyRef.secretName` | If specified, the name of a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                                            | `(none)` |
| `global.otterizeCloud.credentials.clientSecretKeyRef.secretKey`  | If specified, the key for the clientSecret in a pre-created Kubernetes Secret to be used instead of creating a secret with the value of clientSecret.                                        | `(none)` |
| `global.otterizeCloud.apiAddress`                                | Overrides Otterize Cloud default API address.                                                                                                                                                | `(none)` |
| `global.otterizeCloud.apiExtraCAPEMSecret`                       | The name of a secret containing a single `CA.pem` file for an extra root CA used to connect to Otterize Cloud. The secret should be placed in the same namespace as the Otterize deployment. | `(none)` |

## Global parameters

| Key                                           | Description                                                                                                                                                                                                                                                                                                             | Default                              |
|-----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| `global.allowGetAllResources`                 | If defined overrides `allowGetAllResources`.                                                                                                                                                                                                                                                                            |                                      |
| `global.telemetry.enabled`                    | If set to `false`, all anonymous telemetries collection will be disabled                                                                                                                                                                                                                                                | `true`                               |
| `global.telemetry.usage.enabled`              | If set to `false`, collection of anonymous telemetries on product usage will be disabled                                                                                                                                                                                                                                | `true`                               |
| `global.telemetry.errors.enabled`             | If set to `false`, collection of anonymous telemetries on application crashes and errors will be disabled                                                                                                                                                                                                               | `true`                               |
| `global.telemetry.errors.endpointAddress`     | If set, overrides the default endpoint address for anonymous telemetries on application crashes and errors                                                                                                                                                                                                              | `(none)`                             |
| `global.telemetry.errors.stage`               | If set, overrides the default stage for anonymous telemetries on application crashes and errors                                                                                                                                                                                                                         | `(none)`                             |
| `global.telemetry.errors.networkMapperApiKey` | If set, overrides the default API key for anonymous telemetries on application crashes and errors                                                                                                                                                                                                                       | `(none)`                             |
| `global.commonAnnotations`                    | Annotations to add to all deployed objects                                                                                                                                                                                                                                                                              | {}                                   |
| `global.commonLabels`                         | Labels to add to all deployed objects                                                                                                                                                                                                                                                                                   | {}                                   |
| `global.podAnnotations`                       | Annotations to add to all deployed pods                                                                                                                                                                                                                                                                                 | {}                                   |
| `global.podLabels`                            | Labels to add to all deployed pods                                                                                                                                                                                                                                                                                      | {}                                   |
| `global.workloadNameOverrideAnnotationName`   | Which annotation to use (in the [service name resolution algorithm](https://docs.otterize.com/reference/service-identities#kubernetes-service-identity-resolution)) for setting a pod's service name, if not the default. Use this if you already have annotations on your pods that provide the correct workload name. | `intents.otterize.com/workload-name` |
| `global.openshift`                            | Whether to configure and deploy SecurityContextConstraints that allow all components to run with minimal privileges on a default OpenShift installation.                                                                                                                                                                | `false`                              |

## Common parameters

| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default                        |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| `debug`                | Enable debug logs                                                                                                                                                                                                                                                                                                                                                                                                                                             | `false`                        |
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`                         |
