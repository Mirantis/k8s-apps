# ElasticSearch cluster

## Deploy chart
```console
$ helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com
$ helm install mirantisworkloads/elasticsearch
```

## Chart configuration

| Value | Description | Default |
| --- | --- | --- |
| port | Service port | 9200 |
| probeInitialDelaySeconds | Initial delay before probes starting | 90 |
| probeTimeoutSeconds | Timeout for running probes | 5 |
| image.repository | Container image repository | mirantisworkloads/elasticsearch |
| image.tag | Container image tag | 5.2.2 |
| image.pullPolicy | Container pull policy | IfNotPresent |
| master.replicas | Deployment replicas | 2 |
| master.heapSize | JVM option - exact heap size | 1G |
| master.resources.requests.cpu | Container requested cpu | 256m |
| master.resources.requests.memory | Container requested memory | 512Mi |
| master.resources.limits.cpu | Container limited cpu | 512m |
| master.resources.limits.memory | Container limited memory | 1Gi |
| client.replicas | Deployment replicas | 2 |
| client.heapSize | JVM option - exact heap size | 1G |
| client.resources.requests.cpu | Container requested cpu | 256m |
| client.resources.requests.memory | Container requested memory | 512Mi |
| client.resources.limits.cpu | Container limited cpu | 512m |
| client.resources.limits.memory | Container limited memory | 1Gi |
| data.replicas | StatefulSet replicas | 3 |
| data.heapSize | JVM option - exact heap size | 1536m |
| data.resources.requests.cpu | Container requested cpu | 256m |
| data.resources.requests.memory | Container requested memory | 2Gi |
| data.resources.limits.cpu | Container limited cpu | 1 |
| data.resources.limits.memory | Container limited memory | 4Gi |
| data.persistence.type | Mounting volumes type: emptyDir, hostPath or PersistentVolumeClaim | emptyDir |
| data.persistence.storageClass | If type is PersistentVolumeClaim, add persistent storage for it | - |
| data.persistence.volumeSize | Volume size | 10Gi |
| data.persistence.hostPath | Host path for hostPath type of volumes | "" |
