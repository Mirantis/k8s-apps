package client

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	ReplicationController = "ReplicationController"
	Deployment            = "Deployment"
	DaemonSet             = "DaemonSet"
	StatefulSet           = "StatefulSet"
	ReplicaSet            = "ReplicaSet"
	PersistentVolumeClaim = "PersistentVolumeClaim"
	Service               = "Service"
)

func getPods(client *kubernetes.Clientset, namespace string, app string) ([]v1.Pod, error) {
	selector := map[string]string{"app": app}
	list, err := client.Pods(namespace).List(metaV1.ListOptions{
		FieldSelector: fields.Everything().String(),
		LabelSelector: labels.Set(selector).AsSelector().String(),
	})
	return list.Items, err
}

func checkResourcesState(resources map[string][]string, namespace string) (bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return false, err
	}
	client, err := kubernetes.NewForConfig(config)
	pods := []v1.Pod{}
	services := []v1.Service{}
	pvc := []v1.PersistentVolumeClaim{}
	deployments := []*v1beta1.Deployment{}
	for typeRes, res := range resources {
		for _, r := range res {
			switch typeRes {
			case ReplicationController:
				list, err := getPods(client, namespace, r)
				if err != nil {
					return false, err
				}
				pods = append(pods, list...)
			case Deployment:
				deployment, err := client.AppsV1beta1().Deployments(namespace).Get(r, metaV1.GetOptions{})
				if err != nil {
					return false, err
				}
				deployments = append(deployments, deployment)
			case DaemonSet:
				list, err := getPods(client, namespace, r)
				if err != nil {
					return false, err
				}
				pods = append(pods, list...)
			case StatefulSet:
				stSet, err := client.StatefulSets(namespace).Get(r, metaV1.GetOptions{})
				if err != nil {
					return false, err
				}
				list, err := getPods(client, namespace, r)
				if int32(len(list)) < *stSet.Spec.Replicas {
					return false, nil
				}
				if err != nil {
					return false, err
				}
				pods = append(pods, list...)
			case ReplicaSet:
				list, err := getPods(client, namespace, r)
				if err != nil {
					return false, err
				}
				pods = append(pods, list...)
			case PersistentVolumeClaim:
				claim, err := client.PersistentVolumeClaims(namespace).Get(r, metaV1.GetOptions{})
				if err != nil {
					return false, err
				}
				pvc = append(pvc, *claim)
			case Service:
				svc, err := client.Services(namespace).Get(r, metaV1.GetOptions{})
				if err != nil {
					return false, err
				}
				services = append(services, *svc)
			}
		}
	}
	isReady := podsReady(pods) && servicesReady(services) && volumesReady(pvc) && deploymentsReady(deployments)
	return isReady, nil
}

func deploymentsReady(deployments []*v1beta1.Deployment) bool {
	for _, dep := range deployments {
		if dep.Status.ReadyReplicas != int32(*dep.Spec.Replicas) {
			return false
		}
	}
	return true
}

func podsReady(pods []v1.Pod) bool {
	for _, pod := range pods {
		if !v1.IsPodReady(&pod) {
			return false
		}
	}
	return true
}

func servicesReady(svc []v1.Service) bool {
	for _, s := range svc {
		// ExternalName Services are external to cluster so helm shouldn't be checking to see if they're 'ready' (i.e. have an IP Set)
		if s.Spec.Type == v1.ServiceTypeExternalName {
			continue
		}

		// Make sure the service is not explicitly set to "None" before checking the IP
		if s.Spec.ClusterIP != v1.ClusterIPNone && !v1.IsServiceIPSet(&s) {
			return false
		}
		// This checks if the service has a LoadBalancer and that balancer has an Ingress defined
		if s.Spec.Type == v1.ServiceTypeLoadBalancer && s.Status.LoadBalancer.Ingress == nil {
			return false
		}
	}
	return true
}

func volumesReady(vols []v1.PersistentVolumeClaim) bool {
	for _, v := range vols {
		if v.Status.Phase != v1.ClaimBound {
			return false
		}
	}
	return true
}

func getSecretsByNames(names []string, namespace string) (brokerapi.Credential, error) {
	secrets := brokerapi.Credential{}
	config, err := rest.InClusterConfig()
	if err != nil {
		return secrets, err
	}
	client, err := kubernetes.NewForConfig(config)
	for _, name := range names {
		secret, err := client.Secrets(namespace).Get(name, metaV1.GetOptions{})
		if err != nil {
			return secrets, err
		}
		for k, v := range secret.Data {
			secrets[k] = string(v)
		}
	}
	return secrets, nil
}
