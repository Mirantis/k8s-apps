package apply

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mirantis/k8s-apps/helm-apply/kubeutils"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

func (cluster *Cluster) GetHelmClient() (*helm.Client, error) {
	if cluster.helmClient != nil {
		return cluster.helmClient, nil
	}
	for attempts := 60; attempts > 0; attempts-- {
		tunnel, err := cluster.newTunnel()
		if err != nil {
			if strings.Contains(err.Error(), "could not find tiller") {
				err = cluster.initTiller()
				if err != nil {
					return nil, err
				}
				time.Sleep(5 * time.Second)
				continue
			} else if strings.Contains(err.Error(), "could not find a ready tiller pod") {
				time.Sleep(3 * time.Second)
				continue
			} else {
				return nil, err
			}
		}
		tillerHost := fmt.Sprintf("127.0.0.1:%d", tunnel.Local)
		cluster.helmClient = helm.NewClient(helm.Host(tillerHost))
		return cluster.helmClient, nil
	}
	return nil, fmt.Errorf("wasn't able to establish connection with "+
		"tiller for cluster \"%s\"", cluster.Name)
}

func (cluster *Cluster) initTiller() error {
	fmt.Printf("Initializing tiller for cluster \"%s\"...\n", cluster.Name)
	_, client, err := cluster.GetKubeClient()
	if err != nil {
		return err
	}
	namespace := "kube-system"
	if len(cluster.TillerNamespace) > 0 {
		namespace = cluster.TillerNamespace
	}
	err = kubeutils.EnsureNamespace(client, namespace)
	if err != nil {
		return err
	}

	err = kubeutils.CreateServiceAccount(client, namespace)
	if err != nil {
		return err
	}
	err = kubeutils.CreateClusterRoleBinding(client, namespace)
	if err != nil {
		return err
	}

	err = installer.Install(client, &installer.Options{
		Namespace: namespace, ImageSpec: "gcr.io/kubernetes-helm/tiller:v2.6.2",
		ServiceAccount: "tiller-" + namespace})
	if err != nil {
		return err
	}
	return nil
}

func (cluster *Cluster) newTunnel() (*kube.Tunnel, error) {
	config, client, err := cluster.GetKubeClient()
	if err != nil {
		return nil, err
	}
	namespace := "kube-system"
	if len(cluster.TillerNamespace) > 0 {
		namespace = cluster.TillerNamespace
	}
	tunnel, err := portforwarder.New(namespace, client, config)
	if err != nil {
		return nil, err
	}
	return tunnel, nil
}
