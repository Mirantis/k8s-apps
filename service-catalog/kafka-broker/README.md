# Kafka Broker

This repository implmenents a Service Broker which interacts with Kafka to dynamically provision topic and access credentials. This broker should be used with the [Kubernetes Service Broker](https://github.com/kubernetes-incubator/service-catalog) which is currently in development.

# Prerequisites

1. A Kubernetes cluster
2. An installation of [Helm](https://github.com/kubernetes/helm)
3. An installation of Service Catalog (<https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/walkthrough.md>)
4. An installation of Kafka

# Kafka installation

Install the kafka from the `charts` directory in [k8-apps](https://github.com/Mirantis/k8s-apps/):
```bash
$ cd ./charts/kafka
$ helm dep up
$ helm install . --name kafka --namespace kafka-broker
```

# Kafka Broker Installation

Install the kafka broker from the `charts` directory in this repository:

```bash
$ helm install ./charts/kafka-broker --name kafka-broker --namespace kafka-broker --set config.brokers={kafka-kafka-0.kafka-kafka:9092, kafka-kafka-1.kafka-kafka:9092, kafka-kafka-2.kafka-kafka:9092},config.replicationFactor=3,config.partitions=3,config.zookeeperServer=”zk-kafka-0.zk-kafka”
```

# Usage

Namespace createion:

kubectl create namespace test-ns

Broker creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/kafka-broker/examples/kafka-broker.yaml
```
Instance creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/kafka-broker/examples/instance-broker.yaml
```
Binding creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/kafka-broker/examples/binding-broker.yaml
```

After creating these third-party-resources, the Service Catalog and kafka-broker will create a unique topic and access credentials. The topic name and credentials will be written to a secret by the name of "kafka-secret" which is specified in the binding.
