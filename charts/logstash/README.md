LogStash

## Deploy chart
```console
$ helm install .
```

## Chart configuration

| Value | Description | Default |
| --- | --- | --- |
| port | Service port | 5043 |
| replicas | Deployment replicas | 1 |
| elasticsearch.external | If true, logstash uses `host` and `port` values to establish connection with elasticsearch.
                           If false, logstash decides that elasticsearch deployed in same release and uses internal data to connect to elasticsearch. | true |
| elasticsearch.host | Elasticsearch service name to connect | elasticsearch-elasticsearch |
| elasticsearch.port | Elasticsearch service port to connect | 9200 |
| image.repository | Container image repository | 127.0.0.1:31500/logstash |
| image.tag | Container image tag | latest |
| image.pullPolicy | Container pull policy | IfNotPresent |
| heapSize | JVM option - exact heap size | 1536m |
| resources.requests.memory | Container requested memory | 2Gi |
| resources.requests.cpu | Container requested cpu | 250m |
| resources.limits.memory | Container limited memory | 4Gi |
| resources.limits.cpu | Container limited cpu | 1024m |
