# Helm Broker

Universal Helm Service Broker is an implementation of a Service Broker that uses Helm Client for charts provision.
This broker should be used with the [Kubernetes Service Broker](https://github.com/kubernetes-incubator/service-catalog) which is currently in development.

# Prerequisites

1. A Kubernetes cluster 1.6
2. An installation of [Helm](https://github.com/kubernetes/helm) 2.4.2
3. An installation of Service Catalog (<https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/walkthrough.md>)

# Helm Broker Installation

Install the helm broker:
```bash
$ kubectl get pod --namespace=kube-system -o wide
$ export TILLER_IP=172.17.0.4
$ helm install k8s-apps/service-catalog/helm-broker/charts/helm-broker --name helm-broker --namespace helm-broker --set config.tillerHost=$TILLER_IP
```
# Usage

Namespace creation:

```bash
$ kubectl create namespace test-ns
```
Broker creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/helm-broker/examples/helm-broker.yaml
```
Instance creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/helm-broker/examples/helm-instance.yaml
```
Binding creation:

```bash
$ kubectl --context=service-catalog create -f ./charts/helm-broker/examples/helm-binding.yaml
```

After creating these third-party-resources, the Service Catalog and helm-broker will deploy a chart and access credentials. Credentials will be written to a secret by the name of "helm-secret" which is specified in the binding.
