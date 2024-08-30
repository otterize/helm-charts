
# parameters

___

## Global parameters

| Key                                    | Description                                                              | Default         |
|----------------------------------------|--------------------------------------------------------------------------|-----------------|
| `global.spiffe.CASubject`              | The Subject that CA certificates should use (see below)                  |                 |
| `global.spiffe.CASubject.country`      | `Country` value                                                          | `"US"`          |
| `global.spiffe.CASubject.organization` | `Organization` Value                                                     | `"SPIRE"`       |
| `global.spiffe.CASubject.commonName`   | `CommonName` value                                                       | `""`            |
| `global.spiffe.trustDomain`            | The trust domain that this server belongs to                             | `"example.org"` |
| `global.spire.serverServiceName`       | Name of the Kubernetes service that will be created for the SPIRE server |                 |
| `global.commonAnnotations`             | Annotations to add to all deployed objects                               | {}              |
| `global.commonLabels`                  | Labels to add to all deployed objects                                    | {}              |
| `global.podAnnotations`                | Annotations to add to all deployed pods                                  | {}              |
| `global.podLabels`                     | Labels to add to all deployed pods                                       | {}              |

## Common parameters

| Key                                          | Description                                                                                                  | Default             |
|----------------------------------------------|--------------------------------------------------------------------------------------------------------------|---------------------|
| `affinity`                                   | Affinity for pod assignment                                                                                  | `{}`                |
| `autoscaling.enabled`                        | Enable [HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) | `false`             |
| `autoscaling.maxReplicas`                    | HorizontalPodAutoscaler - maxReplicas                                                                        | `100`               |
| `autoscaling.minReplicas`                    | HorizontalPodAutoscaler - minReplicas                                                                        | `1`                 |
| `autoscaling.targetCPUUtilizationPercentage` | HorizontalPodAutoscaler - target CPU utilization percentage                                                  | `80`                |
| `clusterName`                                | Spiffe NodeAttestor Cluster Name                                                                             | `"example-cluster"` |
| `fullnameOverride`                           | String to fully override spire.fullname (from helpers.tpl)                                                   | `""`                |
| `imagePullSecrets`                           | Docker registry secret names as an array                                                                     | `[]`                |
| `nameOverride`                               | String to partially override spire.fullname                                                                  | `""`                |
| `nodeSelector`                               | Node labels for pod assignment                                                                               | `{}`                |
| `podAnnotations`                             | Extra annotations for Spire pods                                                                             | `{}`                |
| `podSecurityContext`                         | Kubernetes [pod SecurityContext](https://jamesdefabia.github.io/docs/user-guide/security-context/)           | `{}`                |
| `replicaCount`                               | Number of replicas                                                                                           | `1`                 |
| `securityContext`                            | Kubernetes [contatiner SecurityContext](https://jamesdefabia.github.io/docs/user-guide/security-context/)    | `{}`                |
| `serviceAccount.annotations`                 | Annotations of the service account                                                                           | `{}`                |
| `serviceAccount.create`                      | Should create service account                                                                                | `true`              |
| `serviceAccount.name`                      | Service account name                                                                                         | `""`                |
| `tolerations`                                | Tolerations for pod assignment                                                                               | `[]`                |

## Agent parameters

| Key                             | Description                       | Default                          |
|---------------------------------|-----------------------------------|----------------------------------|
| `agent.image.pullPolicy`        | Agent image pull policy           | `"IfNotPresent"`                 |
| `agent.image.repository`        | Agent image repository            | `"otterize/spire-agent"` |
| `agent.image.tag`               | Agent image tag                   | `""`                             |
| `agent.logLevel`               | Agent log level                   | `"INFO"`                         |
| `agent.skipKubeletVerification` | Set to `True` if you use Minikube | `false`                          |
| `agent.resources`               | Resources of the container        | `(none)`                         |

## Server parameters

| Key                               | Description                  | Default                            |
|-----------------------------------|------------------------------|------------------------------------|
| `server.dataStorage.accessMode`   | data storage - access mode   | `"ReadWriteOnce"`                  |
| `server.dataStorage.enabled`    | data storage - enabled       | `true`                             |
| `server.dataStorage.size`         | data storage - size          | `"1Gi"`                            |
| `server.dataStorage.storageClass` | data storage - storage class | `nil`                              |
| `server.image.pullPolicy`         | image pull policy            | `"IfNotPresent"`                   |
| `server.image.repository`         | image repository             | `"otterize/spire-server"`  |
| `server.image.tag`                | image tag                    | `""`                               |
| `server.logLevel`                 | log level                    | `"INFO"`                           |
| `server.rootCATTL`                | determine root_ca TTL        | `"26280h"`                         |
| `server.service.type`             | kubernetes service type      | `"ClusterIP"`                      |
| `server.SVIDDefaultTTL`           | determine SVID default TTL   | `"24h"`                            |
| `server.resources`                | Resources of the container   | `(none)`                           |
