Current directory contains all required configs and instructions to run
Twitter analytics demo. Depending on your k8s cluster configuration you could
run services in a stand-alone mode (using config files under `single-node`
directory) or in an ha mode (using config files under `multi-node` directory)

The only things that need to be added to config files are Twitter API
credentials (in `single-node/tweepub.yaml` or `multi-node/tweepub.yaml`):

```console
  twitter:
    appKey: EDITME
    appSecret: EDITME
    tokenKey: EDITME
    tokenSecret: EDITME
```

After that you can install all necessary charts with the following commands.

For single-node:

```console
helm install -n zookeeper-1 ../../charts/zookeeper -f single-node/configs/zookeeper.yaml
helm install -n hdfs-1 ../../charts/hdfs -f single-node/configs/hdfs.yaml

helm install -n kafka-1 ../../charts/kafka -f single-node/configs/kafka.yaml
helm install -n spark-1 ../../charts/spark -f single-node/configs/spark.yaml

helm install -n tweepub-1 ../../charts/tweepub -f single-node/configs/tweepub.yaml
helm install -n tweetics-1 ../../charts/tweetics -f single-node/configs/tweetics.yaml
helm install -n tweeviz-1 ../../charts/tweeviz -f single-node/configs/tweeviz.yaml
```

For multi-node:

```console
helm install -n zookeeper-1 ../../charts/zookeeper -f multi-node/configs/zookeeper.yaml
helm install -n hdfs-1 ../../charts/hdfs -f multi-node/configs/hdfs.yaml

helm install -n kafka-1 ../../charts/kafka -f multi-node/configs/kafka.yaml
helm install -n spark-1 ../../charts/spark -f multi-node/configs/spark.yaml

helm install -n tweepub-1 ../../charts/tweepub -f multi-node/configs/tweepub.yaml
helm install -n tweetics-1 ../../charts/tweetics -f multi-node/configs/tweetics.yaml
helm install -n tweeviz-1 ../../charts/tweeviz -f multi-node/configs/tweeviz.yaml
```

After that you can access tweeviz endpoint to see a tag cloud.