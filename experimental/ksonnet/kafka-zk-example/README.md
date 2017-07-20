# Kafka + ZK app

Standard kafka + zk app with features, implemented in same helm chart.

# How to install

1. Add new lib path to KUBECFG_JPATH:

   ```console
   export KUBECFG_JPATH=$KUBECFG_JPATH:<absolute path>/k8s-apps/experimental/ksonnet
   ```

2. Show yaml template with command:

   ```console
   kubecfg show -o yaml kafka-zk-example/kafka-zk-ps.jsonnet
   ```

3. Create app with kubecfg:

   ```console
   kubecfg update kafka-zk-example/kafka-zk-ps.jsonnet
   ```

To delete app, enter command:

```console
kubecfg delete kafka-zk-example/kafka-zk-ps.jsonnet
```
