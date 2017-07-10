# Twitter Stats

Current directory contains all required configs and instructions to run
Twitter Stats demo. This demo includes following services:

* ZooKeeper
* Kafka
* Spark
* HDFS (or Cassandra)

and shows following workloads flow:

* `tweepub` pulls new tweets from the Twitter Streaming API and puts them into Kafka
* `tweetics` is a Spark Streaming job that takes tweets from Kafka, process them in micro-batches (5 seconds long)
  and puts processed blocks into HDFS or Cassandra based on selection
* `tweeviz` displays word cloud for most popular tweets in some location

# How to

1. Make sure to have kubectl configured to the k8s cluster that you'd like to use.

1. Switch to `scenarios/twitter-stats` directory. There is a bash script `./twitter-stats.sh`
   that runs helm with needed variables, helps to wait when everything up and running and helps
   to cleanup env when you're done.

1. First of all you need to have app created in Twitter to be able to access Twitter Streaming API. You can do it on 
  https://apps.twitter.com/ 
  After creating application, you should create token, it should read-write.
  You should set environment variables for this app and token credentials.

```bash
export TS_APP_KEY=ABCDERFG
export TS_APP_SECRET=ABCDERFG
export TS_TOKEN_KEY=123456-ABCDERFG
export TS_TOKEN_SECRET=ABCDERFG
```

1. You have a few more configs that you can set, all using environment variables with prefix
  `TS_` and there are a few available commands. You can run `./twitter-stats.sh` to see
  help message that (hopefully) covers how to use this tool.
  
```
Commands (multiple commads could be used at the same time):
    up - deploy or upgrade Twitter Stats pipeline with all dependencies, not waiting for completion
         (it's safe to re-run up many times to apply changes to the configs)
	test - wait and test that Twitter Stats demo is working and displaying some stats
	down - destroy Twitter Stats deployment (requires only NAMESPACE to be set)

Parameters (could be set using environment variables):
	TS_NAME - name of the Twitter Stats deployment
              change it to do multiple deployments even in single namespace (used in pod names, etc.)
 		      use 'rand' to generate random name
 	TS_NAMESPACE - K8s namespace that will be used for deployment
 	TS_DELETE_NS - Delete or not K8s namespace when 'down' called
 	TS_APP_KEY, TS_APP_SECRET, TS_TOKEN_KEY, TS_TOKEN_SECRET - Twitter API credentials (must be read+write)
 	TS_MODE - single or multi-node deployment
 	TS_STORAGE - hdfs or cassandra for storing processed data
 	TS_RETRIES - number of retries for 'test' command
 	TS_RETRY_INTERVAL - amount of time in seconds to sleep between retries
 	TS_KUBECTL_CMD - path to kubectl binary to run
 	TS_HELM_CMD - path to helm binary to run
```

1. Additionally you can edit Helm config files in the `single-node/configs` or `multi-node/configs`
   directory (depends on selected mode `TS_MODE`).

1. To run demo and await fully working demo just run `./twitter-stats.sh up test`. Then you can
   rerun same command multiple times to re-deploy updated configs. For example, you can edit
   file `multi-node/configs/kafka.yaml` (or `single-`) and increase number of kafka replicas from 3
   to 4. To see all pods in the `demo` namespace, you can run `kubectl -n demo get pods -o wide`.
   
1. When you're done showing demo, you can run `./twitter-stats.sh down` to cleanup env (deletes namespace
   by default).

# Sample run

Here is the sample run of `./twitter-stats.sh up test` (you'll see in the end of the output
link to the Twitter Stats visualization like `Link: http://104.196.56.74:30103`):

```
~/w/k/k/s/twitter-stats ➜ ./twitter-stats.sh up test
[2017-07-09 22:55:56] twitter-stats [main] | --------------------------------------------------
[2017-07-09 22:55:56] twitter-stats [main] | Welcome to the Twitter Stats!
[2017-07-09 22:55:56] twitter-stats [main] | --------------------------------------------------
[2017-07-09 22:55:56] twitter-stats [main] | Legend (available commands):
[2017-07-09 22:55:56] twitter-stats [main] | 	up - deploy or upgrade Twitter Stats pipeline with all dependencies, not waiting for completion
[2017-07-09 22:55:56] twitter-stats [main] | 		(it's safe to re-run up many times to apply changes to the configs)
[2017-07-09 22:55:56] twitter-stats [main] | 	wait - waiting for Twitter Stats ready and returns stats (requires only NAMESPACE to be set)
[2017-07-09 22:55:56] twitter-stats [main] | 	down - destroy Twitter Stats deployment (requires only NAMESPACE to be set)
[2017-07-09 22:55:56] twitter-stats [main] | 	test - deploy, wait and destroy Twitter Stats
[2017-07-09 22:55:56] twitter-stats [main] |
[2017-07-09 22:55:56] twitter-stats [main] | Parameters (could be set using environment variables):
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_NAME - name of the Twitter Stats deployment
[2017-07-09 22:55:56] twitter-stats [main] | 		change it to do multiple deployments even in single namespace (used in pod names, etc.)
[2017-07-09 22:55:56] twitter-stats [main] | 		use 'rand' to generate random name
[2017-07-09 22:55:56] twitter-stats [main] | 		current: demo (default: demo)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_NAMESPACE - K8s namespace that will be used for deployment
[2017-07-09 22:55:56] twitter-stats [main] | 		current: demo (default: same as TS_NAME)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_DELETE_NS - Delete or not K8s namespace when 'down' called
[2017-07-09 22:55:56] twitter-stats [main] | 		current: yes (default: yes)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_APP_KEY, TS_APP_SECRET, TS_TOKEN_KEY, TS_TOKEN_SECRET - Twitter API credentials (must be read+write)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_MODE - single or multi-node deployment
[2017-07-09 22:55:56] twitter-stats [main] | 		current: multi (default: multi)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_STORAGE - hdfs or cassandra for storing processed data
[2017-07-09 22:55:56] twitter-stats [main] | 		current: hdfs (default: hdfs)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_RETRIES - number of retries for 'test' command
[2017-07-09 22:55:56] twitter-stats [main] | 		current: 60 (default: 60)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_RETRY_INTERVAL - amount of time in seconds to sleep between retries
[2017-07-09 22:55:56] twitter-stats [main] | 		current: 10 (default: 10)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_KUBECTL_CMD - path to kubectl binary to run
[2017-07-09 22:55:56] twitter-stats [main] | 		current: kubectl (default: kubectl)
[2017-07-09 22:55:56] twitter-stats [main] | 	TS_HELM_CMD - path to helm binary to run
[2017-07-09 22:55:56] twitter-stats [main] | 		current: helm (default: helm)
[2017-07-09 22:55:56] twitter-stats [main] | --------------------------------------------------
[2017-07-09 22:55:56] twitter-stats [main] | Following commands will be executed: up test
[2017-07-09 22:55:56] twitter-stats [main] | --------------------------------------------------
[2017-07-09 22:55:56] twitter-stats [  up] | --------------------------------------------------
[2017-07-09 22:55:56] twitter-stats [  up] | Deploying or Upgrading Twitter Stats services
[2017-07-09 22:55:56] twitter-stats [  up] | --------------------------------------------------
namespace "demo" created
[2017-07-09 22:55:58] twitter-stats [  up] | Calculated configs and Helm logs: /var/folders/xv/gjzz8x0s0xz4781p240h090m0000gn/T/tmp.cBJ6nwF0
[2017-07-09 22:55:59] twitter-stats [  up] | Setting up Helm
[2017-07-09 22:56:03] twitter-stats [  up] | Deploying or upgrading chart zookeeper with release ts-demo-zookeeper
Release "ts-demo-zookeeper" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:09 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> v1/ConfigMap
NAME                  DATA  AGE
zk-ts-demo-zookeeper  3     1s

==> v1/Service
NAME                  CLUSTER-IP  EXTERNAL-IP  PORT(S)                     AGE
zk-ts-demo-zookeeper  None        <none>       2888/TCP,3888/TCP,2181/TCP  1s

==> apps/v1beta1/StatefulSet
NAME                  DESIRED  CURRENT  AGE
zk-ts-demo-zookeeper  3        0        1s

==> policy/v1beta1/PodDisruptionBudget
NAME                  MIN-AVAILABLE  ALLOWED-DISRUPTIONS  AGE
zk-ts-demo-zookeeper  2              0                    1s


NOTES:
ZooKeeper chart has been deployed.

Internal URL:
    zookeeper: zk-ts-demo-zookeeper-0.zk-ts-demo-zookeeper:2181,zk-ts-demo-zookeeper-1.zk-ts-demo-zookeeper:2181,zk-ts-demo-zookeeper-2.zk-ts-demo-zookeeper:2181

[2017-07-09 22:56:08] twitter-stats [  up] | Deploying or upgrading chart hdfs with release ts-demo-hdfs
Release "ts-demo-hdfs" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:15 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> v1/ConfigMap
NAME                       DATA  AGE
hdfs-configs-ts-demo-hdfs  1     1s

==> v1/Service
NAME                        CLUSTER-IP    EXTERNAL-IP  PORT(S)          AGE
hdfs-namenode-ts-demo-hdfs  None          <none>       8020/TCP         1s
hdfs-datanode-ts-demo-hdfs  None          <none>       50075/TCP        1s
hdfs-ui-ts-demo-hdfs        10.23.248.74  <nodes>      50070:31606/TCP  1s

==> apps/v1beta1/StatefulSet
NAME                        DESIRED  CURRENT  AGE
hdfs-datanode-ts-demo-hdfs  3        0        1s
hdfs-namenode-ts-demo-hdfs  1        0        1s


NOTES:
HDFS chart has been deployed.

Internal URL:
    namenode: hdfs-namenode-ts-demo-hdfs-0.hdfs-namenode-ts-demo-hdfs:8020
    hdfs-ui: hdfs-ui-ts-demo-hdfs:50070

External URL:
Get the HDFS UI URL to visit by running these commands in the same shell:
    export NODE_PORT=$(kubectl get --namespace demo -o jsonpath="{.spec.ports[0].nodePort}" services hdfs-ui-ts-demo-hdfs)export NODE_IP=$(kubectl get nodes --namespace demo -o jsonpath="{.items[0].status.addresses[0].address}")
    echo http://$NODE_IP:$NODE_PORT


[2017-07-09 22:56:14] twitter-stats [  up] | Deploying or upgrading chart kafka with release ts-demo-kafka
Release "ts-demo-kafka" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:21 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> v1/ConfigMap
NAME                    DATA  AGE
kafka-fb-ts-demo-kafka  1     1s

==> v1/Service
NAME                 CLUSTER-IP  EXTERNAL-IP  PORT(S)   AGE
kafka-ts-demo-kafka  None        <none>       9092/TCP  1s

==> apps/v1beta1/StatefulSet
NAME                 DESIRED  CURRENT  AGE
kafka-ts-demo-kafka  3        0        1s


NOTES:
Kafka chart has been deployed.

Internal URL:
    kafka: kafka-ts-demo-kafka-0.kafka-ts-demo-kafka:9092,kafka-ts-demo-kafka-1.kafka-ts-demo-kafka:9092,kafka-ts-demo-kafka-2.kafka-ts-demo-kafka:9092

[2017-07-09 22:56:21] twitter-stats [  up] | Deploying or upgrading chart spark with release ts-demo-spark
Release "ts-demo-spark" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:28 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> extensions/v1beta1/Deployment
NAME                        DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
zeppelin-ts-demo-spark      1        1        1           0          2s
spark-worker-ts-demo-spark  3        3        3           0          2s

==> apps/v1beta1/StatefulSet
NAME                        DESIRED  CURRENT  AGE
spark-master-ts-demo-spark  1        0        2s

==> v1/ConfigMap
NAME                       DATA  AGE
spark-conf-ts-demo-spark   1     2s
zeppelin-fb-ts-demo-spark  1     2s
spark-fb-ts-demo-spark     1     2s

==> v1/Service
NAME                            CLUSTER-IP     EXTERNAL-IP  PORT(S)                                       AGE
spark-master-ext-ts-demo-spark  10.23.240.175  <nodes>      7077:31715/TCP,8080:31649/TCP,6066:32607/TCP  2s
spark-master-ts-demo-spark      None           <none>       7077/TCP                                      2s
zeppelin-ts-demo-spark          10.23.247.13   <nodes>      8080:32487/TCP                                2s


NOTES:
Spark chart has been deployed.

Internal URL:
    spark: spark-master-ts-demo-spark-0.spark-master-ts-demo-spark:7077
    zeppelin: zeppelin-ts-demo-spark:8080

External URL:
To get the Spark and Zeppelin URL to visit by running these commands in the same shell:

    export NODE_IP=$(kubectl get nodes --namespace demo -o jsonpath="{.items[0].status.addresses[0].address}")
    export SPARK_NODE_PORT=$(kubectl get --namespace demo -o jsonpath="{.spec.ports[1].nodePort}" services spark-master-ts-demo-spark)
    export ZEPPELIN_NODE_PORT=$(kubectl get --namespace demo -o jsonpath="{.spec.ports[0].nodePort}" services zeppelin-ts-demo-spark)
    echo http://$NODE_IP:$SPARK_NODE_PORT
    echo http://$NODE_IP:$ZEPPELIN_NODE_PORT
[2017-07-09 22:56:31] twitter-stats [  up] | Deploying or upgrading chart tweepub with release ts-demo-tweepub
Release "ts-demo-tweepub" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:37 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> extensions/v1beta1/Deployment
NAME                     DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
tweepub-ts-demo-tweepub  1        1        1           0          1s


NOTES:
TweePub chart has been deployed.

[2017-07-09 22:56:38] twitter-stats [  up] | Deploying or upgrading chart tweetics with release ts-demo-tweetics
Release "ts-demo-tweetics" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:45 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> extensions/v1beta1/Deployment
NAME                       DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
tweetics-ts-demo-tweetics  1        1        1           0          1s


NOTES:
Tweetics job has been started.

[2017-07-09 22:56:44] twitter-stats [  up] | Deploying or upgrading chart tweeviz with release ts-demo-tweeviz
Release "ts-demo-tweeviz" has been upgraded. Happy Helming!
LAST DEPLOYED: Sun Jul  9 22:56:51 2017
NAMESPACE: demo
STATUS: DEPLOYED

RESOURCES:
==> v1/Service
NAME                     CLUSTER-IP     EXTERNAL-IP  PORT(S)         AGE
tweeviz-ts-demo-tweeviz  10.23.243.206  <nodes>      8589:30103/TCP  1s

==> extensions/v1beta1/Deployment
NAME                     DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
tweeviz-ts-demo-tweeviz  1        1        1           0          1s


NOTES:
TweeViz chart has been deployed.

[2017-07-09 22:56:53] twitter-stats [  up] | All Twitter Stats services deployed or upgraded (async)
[2017-07-09 22:56:53] twitter-stats [test] | --------------------------------------------------
[2017-07-09 22:56:53] twitter-stats [test] | Wait Twitter Stats ready and serve stats
[2017-07-09 22:56:53] twitter-stats [test] | --------------------------------------------------
[2017-07-09 22:56:53] twitter-stats [test] | Will retry 60 times with 10 seconds intervals
[2017-07-09 22:56:56] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (0/60)
[2017-07-09 22:57:08] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (1/60)
[2017-07-09 22:57:19] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (2/60)
[2017-07-09 22:57:31] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (3/60)
[2017-07-09 22:57:44] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (4/60)
[2017-07-09 22:57:56] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (5/60)
[2017-07-09 22:58:08] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (6/60)
[2017-07-09 22:58:20] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (7/60)
[2017-07-09 22:58:32] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (8/60)
[2017-07-09 22:58:44] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (9/60)
[2017-07-09 22:58:58] twitter-stats [test] | Twitter Stats isn't ready or not serving stats (yet), sleeping for 10 seconds (10/60)
[2017-07-09 22:59:10] twitter-stats [test] | --------------------------------------------------
[2017-07-09 22:59:10] twitter-stats [test] | Twitter Stats ready and serve stats successfully
[2017-07-09 22:59:10] twitter-stats [test] | --------------------------------------------------
[2017-07-09 22:59:10] twitter-stats [test] | Link: http://104.196.56.74:30103
[2017-07-09 22:59:10] twitter-stats [test] |
[2017-07-09 22:59:10] twitter-stats [test] | Successfully finished
[2017-07-09 22:59:10] twitter-stats [exit] |
[2017-07-09 22:59:10] twitter-stats [exit] | Total execution time: 194 seconds
```
