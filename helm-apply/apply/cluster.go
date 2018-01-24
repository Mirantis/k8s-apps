package apply

import (
	"fmt"
	"strconv"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/helm"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster struct {
	Path            string `yaml:"path"`
	Context         string `yaml:"context"`
	ExternalIP      string `yaml:"externalIP"`
	Namespace       string `yaml:"namespace"`
	TillerNamespace string `yaml:"tillerNamespace"`
	Name            string `yaml:"-"`
	helmClient      *helm.Client
	kubeClient      kubernetes.Interface
	kubeConfig      *rest.Config
}

func (cluster *Cluster) getConfig() clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(cluster.Path) > 0 {
		rules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: cluster.Path}
	}

	overrides := &clientcmd.ConfigOverrides{}
	if len(cluster.Context) > 0 {
		overrides.CurrentContext = cluster.Context
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

func (cluster *Cluster) GetKubeClient() (*rest.Config, kubernetes.Interface, error) {
	if cluster.kubeClient != nil && cluster.kubeConfig != nil {
		return cluster.kubeConfig, cluster.kubeClient, nil
	}
	loadingConfig := cluster.getConfig()
	conf, err := loadingConfig.ClientConfig()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not get Kubernetes config for cluster \"%s\": %s", cluster.Name, err)
	}
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not get Kubernetes client for cluster %s: %s", cluster.Name, err)
	}
	cluster.kubeConfig = conf
	cluster.kubeClient = client
	return conf, client, nil
}

func (cluster *Cluster) GetDefaultNamespace() string {
	if len(cluster.Namespace) > 0 {
		return cluster.Namespace
	}
	loadingConfig := cluster.getConfig()
	if ns, _, err := loadingConfig.Namespace(); err == nil {
		return ns
	}
	return "default"
}

func (cluster *Cluster) serviceExternalAddress(serviceName, portName, namespace string) (string, error) {
	_, client, err := cluster.GetKubeClient()
	if err != nil {
		return "", err
	}
	svc, err := client.CoreV1().Services(namespace).Get(serviceName, v1.GetOptions{})
	if err != nil {
		return "", err
	}
	switch svc.Spec.Type {
	case apiv1.ServiceTypeNodePort:
		if cluster.ExternalIP == "" {
			return "", fmt.Errorf("externalIP for cluster \"%s\" is not provided", cluster.Name)
		}
		for _, port := range svc.Spec.Ports {
			intPort, err := strconv.Atoi(portName)
			var int32Port int32
			if err == nil {
				int32Port = int32(intPort)
			}
			if port.Name == portName || port.Port == int32Port {
				return fmt.Sprintf("%s:%d", cluster.ExternalIP, port.NodePort), nil
			}
		}
		return "", fmt.Errorf("port \"%s\" for service \"%s\" was not found", serviceName, portName)
	}
	return "", fmt.Errorf("cannot get external address for \"%s\" service", serviceName)
}
