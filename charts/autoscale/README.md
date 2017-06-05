# AutoScale

## Deploy chart
```console
$ helm install .
```

AutoScale works with K8s StatefulSets and Deployments.
The following annotations are supported:

autoscale/minReplicas - object will not be scaled down to the lesser number of replicas
                        than defined here
autoscale/maxReplicas - object will not be scaled up to the larger number of replicas
                        than defined here
autoscale/up          - the prometheus query that triggers scale up if conditions are met
autoscale/down        - the prometheus query that triggers scale down if conditions are met
