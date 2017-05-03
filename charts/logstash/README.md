LogStash

## Deploy chart
```console
$ helm install .
```

## Chart configuration

| Value | Description | Default |
| --- | --- | --- |
| component | Service name | logstash |
| port | Service port | 5043 |
| protocol | Protocol to connect to service | TCP |
| replicas | Deployment replicas | 1 |
| elasticsearch.service | Elasticsearch service name to connect | elasticsearch-elasticsearch |
| elasticsearch.port | Elasticsearch service port to connect | 9200 |
| image.repository | Container image repository | 127.0.0.1:31500/logstash |
| image.tag | Container image tag | latest |
| image.pullPolicy | Container pull policy | Always |
| heapSize | JVM option - exact heap size | 1536m |
| resources.requests.memory | Container requested memory | 2Gi |
| resources.requests.cpu | Container requested cpu | 250m |
| resources.limits.memory | Container limited memory | 4Gi |
| resources.limits.cpu | Container limited cpu | 1024m |
