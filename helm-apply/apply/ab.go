package apply

import (
	"strings"

	"fmt"

	proxyconfig "istio.io/api/proxy/v1/config"
	"istio.io/istio/pilot/adapter/config/crd"
	"istio.io/istio/pilot/model"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type AB struct {
	Name          string          `yaml:"-"`
	ClusterName   string          `yaml:"cluster"`
	Cluster       *Cluster        `yaml:"-"`
	Namespace     string          `yaml:"namespace"`
	IngressDomain string          `yaml:"ingressDomain"`
	Route         []*ReleaseRoute `yaml:"route"`
}

type ReleaseRoute struct {
	ReleaseName string   `yaml:"release"`
	Release     *Release `yaml:"-"`
	Weight      int32    `yaml:"weight"`
}

func (ab AB) IstioClient() (client *crd.Client, err error) {
	istioConfigTypes := model.ConfigDescriptor{
		model.RouteRule,
	}
	client, err = crd.NewClient(ab.Cluster.Path, istioConfigTypes, "")
	if err != nil {
		return nil, err
	}
	return
}

func (ab AB) getNamespace() (namespace string) {
	namespace = ab.Namespace
	if namespace == "" {
		namespace = ab.Cluster.GetDefaultNamespace()
	}
	return
}

func (ab AB) Create() error {
	fmt.Printf("Creating %q route rule...\n", ab.Name)
	kube := ab.Cluster.kubeClient
	istio, err := ab.IstioClient()
	if err != nil {
		return err
	}
	routeSpecMap := map[string]*proxyconfig.RouteRule{}
	namespace := ab.getNamespace()
	for i, route := range ab.Route {
		selector := labels.Set{"release": route.ReleaseName}.AsSelector().String()
		options := meta.ListOptions{LabelSelector: selector}
		services, err := kube.CoreV1().Services(namespace).List(options)
		if err != nil {
			return err
		}
		for _, service := range services.Items {
			if service.Spec.ClusterIP == v1.ClusterIPNone {
				continue
			}
			destName := strings.Replace(service.Name, route.ReleaseName, ab.Name, 1)
			if i == 0 {
				routeSpec := &proxyconfig.RouteRule{
					Destination: &proxyconfig.IstioService{Name: destName, Namespace: namespace},
				}
				routeSpecMap[destName] = routeSpec
				newPorts := service.Spec.Ports[:0]

				for _, port := range service.Spec.Ports {
					port.NodePort = 0
					newPorts = append(newPorts, port)
				}
				destService := v1.Service{
					ObjectMeta: meta.ObjectMeta{
						Name:      destName,
						Namespace: namespace,
					},
					Spec: v1.ServiceSpec{
						Ports: newPorts,
						Type: "ClusterIP",
						Selector: map[string]string{"chartName": route.Release.Chart.Name},
					},
				}

				_, err = kube.CoreV1().Services(namespace).Create(&destService)
				if err != nil && !errors.IsAlreadyExists(err)  {
					return err
				}
				ingSelector := labels.Set{"app": service.Labels["app"]}.AsSelector().String()
				ingOptions := meta.ListOptions{LabelSelector: ingSelector}
				ingresses, err := kube.ExtensionsV1beta1().Ingresses(namespace).List(ingOptions)
				if err != nil {
					return err
				}
				if len(ingresses.Items) == 0 {
					continue
				}
				destIngress := v1beta1.Ingress{
					ObjectMeta: meta.ObjectMeta{
						Name:      destName,
						Namespace: namespace,
						Annotations: map[string]string{"kubernetes.io/ingress.class": "istio"},
					},
					Spec: ingresses.Items[0].Spec,
				}
				destIngress.Spec.Rules[0].Host = ab.IngressDomain
				destIngress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName = destName
				_, err = kube.ExtensionsV1beta1().Ingresses(namespace).Create(&destIngress)
				if err != nil && !errors.IsAlreadyExists(err) {
					return err
				}
			}
			destWeight := &proxyconfig.DestinationWeight{
				Labels: service.Spec.Selector,
				Weight: route.Weight,
			}
			routeSpecMap[destName].Route = append(routeSpecMap[destName].Route, destWeight)
		}
	}
	for name, routeSpec := range routeSpecMap {
		config := model.Config{
			ConfigMeta: model.ConfigMeta{Type: model.RouteRule.Type, Name: name, Namespace: namespace},
			Spec:       routeSpec,
		}
		istio.Delete(model.RouteRule.Type, name, namespace)
		_, err = istio.Create(config)
		if err != nil {
			return err
		}
	}
	return nil
}
