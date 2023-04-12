# Parameters

## Mapper parameters
| Key                            | Description                          | Default                        |
|--------------------------------|--------------------------------------|--------------------------------|
| `mapper.image.repository`      | Mapper image repository.             | `otterize`                     |
| `mapper.image.image`           | Mapper image.                        | `network-mapper`               |
| `mapper.image.tag`             | Mapper image tag.                    | `latest`                       |
| `mapper.pullPolicy`            | Mapper pull policy.                  | `(none)`                       |
| `mapper.pullSecrets`           | Mapper pull secrets.                 | `(none)`                       |
| `mapper.resources`             | Resources override.                  | `(none)`                       |
| `mapper.uploadIntervalSeconds` | Interval for uploading data to cloud | `60`                           |

## Sniffer parameters
| Key                        | Description                | Default                  |
|----------------------------|----------------------------|--------------------------|
| `sniffer.enable`           | Enable sniffer deployment. | `true`                   |
| `sniffer.image.repository` | Sniffer image repository.  | `otterize`               |
| `sniffer.image.image`      | Sniffer image.             | `network-mapper-sniffer` |
| `sniffer.image.tag`        | Sniffer image tag.         | `latest`                 |
| `sniffer.pullPolicy`       | Sniffer pull policy.       | `(none)`                 |
| `sniffer.pullSecrets`      | Sniffer pull secrets.      | `(none)`                 |
| `sniffer.resources`        | Resources override.        | `(none)`                 |   


## Kafka watcher parameters
| Key                             | Description                                                 | Default                        |
|---------------------------------|-------------------------------------------------------------|--------------------------------|
| `kafkawatcher.enable`           | Enable Kafka watcher deployment (experimental).             | `false`                        |
| `kafkawatcher.image.repository` | Kafka watcher image repository.                             | `otterize`                     |
| `kafkawatcher.image.image`      | Kafka watcher image.                                        | `network-mapper-kafka-watcher` |
| `kafkawatcher.image.tag`        | Kafka watcher image tag.                                    | `latest`                       |
| `kafkawatcher.pullPolicy`       | Kafka watcher pull policy.                                  | `(none)`                       |
| `kafkawatcher.pullSecrets`      | Kafka watcher pull secrets.                                 | `(none)`                       |
| `kafkawatcher.resources`        | Resources override.                                         | `(none)`                       |
| `kafkawatcher.kafkaServers`     | Kafka servers to watch, specified as `pod.namespace` items. | `(none)`                       |

## Istio watcher parameters
| Key                             | Description                                     | Default                        |
|---------------------------------|-------------------------------------------------|--------------------------------|
| `istiowatcher.enable`           | Enable Istio watcher deployment (experimental). | `false`                        |
| `istiowatcher.image.repository` | Istio watcher image repository.                 | `otterize`                     |
| `istiowatcher.image.image`      | Istio watcher image.                            | `network-mapper-istio-watcher` |
| `istiowatcher.image.tag`        | Istio watcher image tag.                        | `latest`                       |
| `istiowatcher.pullPolicy`       | Istio watcher pull policy.                      | `(none)`                       |
| `istiowatcher.pullSecrets`      | Istio watcher pull secrets.                     | `(none)`                       |
| `istiowatcher.resources`        | Resources override.                             | `(none)`                       |

## Cloud parameters
| Key                                             | Description                                                                                                                                                                                  | Default  |
|-------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `global.otterizeCloud.credentials.clientId`     | Client ID for connecting to Otterize Cloud.                                                                                                                                                  | `(none)` |
| `global.otterizeCloud.credentials.clientSecret` | Client secret for connecting to Otterize Cloud.                                                                                                                                              | `(none)` |
| `global.otterizeCloud.apiAddress`               | Overrides Otterize Cloud default API address.                                                                                                                                                | `(none)` |
| `global.otterizeCloud.apiExtraCAPEMSecret`      | The name of a secret containing a single `CA.pem` file for an extra root CA used to connect to Otterize Cloud. The secret should be placed in the same namespace as the Otterize deployment. | `(none)` |

## Global parameters
| Key                              | Description                                                                                                                                 | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `global.allowGetAllResources`    | If defined overrides `allowGetAllResources`.                                                                                                |         |

## Common parameters
| Key                    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Default                        |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| `debug`                | Enable debug logs                                                                                                                                                                                                                                                                                                                                                                                                                                             | `false`                        |
| `allowGetAllResources` | Gives get, list and watch permission to watch on all resources. This is used to resolve service names when pods have owners that are custom resources. When disabled, a limited set of permissions is used that only allows access to built-in Kubernetes resources that deploy Pods and Pods themselves - Deployments, StatefulSets, DaemonSets, ReplicaSets and Services. Resolving may not be able to complete if the owning resource is not one of those. | `true`                         |
