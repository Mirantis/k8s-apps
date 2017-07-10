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
