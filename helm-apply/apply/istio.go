package apply

import (
	"fmt"

	proxyconfig "istio.io/api/proxy/v1/config"
	"istio.io/istio/pilot/adapter/config/crd"
	"istio.io/istio/pilot/model"
)

type Route struct {
	Name        string        `yaml:"name"`
	ClusterName string        `yaml:"cluster"`
	Cluster     *Cluster      `yaml:"-"`
	Namespace   string        `yaml:"namespace"`
	Spec        *NewRouteRule `yaml:"spec"`
}

type NewRouteRule struct {
	proxyconfig.RouteRule `yaml:",inline"`
	Route []*NewDestinationWeight
}

type NewDestinationWeight struct {
	proxyconfig.DestinationWeight `yaml:",inline"`
	//Destination *NewIstioService `protobuf:"bytes,1,opt,name=destination" json:"destination,omitempty"`
}

type NewIstioService struct {
	proxyconfig.IstioService `yaml:",inline"`
	//Release string `yaml:"release"`
}

func (r Route) getNamespace() (namespace string) {
	namespace = r.Namespace
	if namespace == "" {
		namespace = r.Cluster.GetDefaultNamespace()
	}
	return
}

func (r Route) IstioClient() (client *crd.Client, err error) {
	istioConfigTypes := model.ConfigDescriptor{
		model.RouteRule,
	}
	client, err = crd.NewClient(r.Cluster.Path, istioConfigTypes, "")
	if err != nil {
		return nil, err
	}
	return
}

func (r Route) Create() error {
	client, err := r.IstioClient()
	if err != nil {
		return err
	}
	namespace := r.getNamespace()
	config := model.Config{
		ConfigMeta: model.ConfigMeta{Type: model.RouteRule.Type, Name: r.Name, Namespace: namespace},
		Spec:       r.Spec,
	}
	client.Delete(model.RouteRule.Type, r.Name, namespace)
	fmt.Printf("Creating route %q...\n", r.Name)
	_, err = client.Create(config)
	if err != nil {
		return err
	}
	return nil
}
