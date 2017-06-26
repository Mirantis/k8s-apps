package controller

import (
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"helm-broker/client"
	"helm-broker/utils"
	"k8s.io/helm/pkg/helm"
	"log"
)

type helmController struct {
	helmClient helm.Client
	chartUrl   string
	tillerHost string
}

type Controller interface {
	Catalog() (*brokerapi.Catalog, error)

	GetServiceInstance(id string) (*brokerapi.LastOperationResponse, error)
	CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (*brokerapi.CreateServiceInstanceResponse, error)
	RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error)

	Bind(instanceID string, bindingID string, req *brokerapi.BindingRequest) (*brokerapi.CreateServiceBindingResponse, error)
	UnBind(instanceID string, bindingID string) error
}

// CreateController returns a Helm Broker Controller
func CreateController(c Config) (Controller, error) {

	helmClient := helm.NewClient(helm.Host(c.TillerHost))

	return &helmController{
		helmClient: *helmClient,
		chartUrl:   c.ChartUrl,
		tillerHost: c.TillerHost,
	}, nil
}

// Catalog returns the Helm Broker catalog entries
func (c *helmController) Catalog() (*brokerapi.Catalog, error) {
	err := utils.DownloadIndex(c.chartUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	index, err := utils.ParseIndex()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var services []*brokerapi.Service
	for _, val := range index.Entries {
		for _, v := range val {
			service := &brokerapi.Service{
				Name:        v.Name + v.Version,
				ID:          uuid.NewV4().String(),
				Description: v.Description,
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          uuid.NewV4().String(),
						Description: v.Description,
						Free:        true,
					},
				},
				Bindable: true,
			}
			services = append(services, service)
		}
	}
	return &brokerapi.Catalog{
		Services: services,
	}, nil
}

// CreateServiceInstance
func (c *helmController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (
	*brokerapi.CreateServiceInstanceResponse, error) {
	err := validateInstanceRequest(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	name := req.Parameters["name"].(string)
	version := req.Parameters["version"].(string)
	namespace := req.Parameters["namespace"].(string)
	values, _ := yaml.Marshal(req.Parameters["values"])

	chartPath, err := utils.DownloadChart(name, version)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = client.Install(c.helmClient, chartPath, id, namespace, values)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &brokerapi.CreateServiceInstanceResponse{
		DashboardURL: "",
		Operation:    "chart installation",
	}, nil
}

// GetServiceInstance
func (c *helmController) GetServiceInstance(id string) (*brokerapi.LastOperationResponse, error) {
	isReady, err := client.IsResourcesReady(c.helmClient, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if isReady {
		log.Printf("Release %s was successfully deployed", id)
		return &brokerapi.LastOperationResponse{
			State:       brokerapi.StateSucceeded,
			Description: "Chart was successfully deployed",
		}, nil
	}
	return &brokerapi.LastOperationResponse{
		State:       brokerapi.StateInProgress,
		Description: "Chart deployment is in progress",
	}, nil
}

// RemoveServiceInstance
func (c *helmController) RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error) {
	err := client.Delete(c.helmClient, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

// Bind
func (c *helmController) Bind(instanceID string, bindingID string, req *brokerapi.BindingRequest) (
	*brokerapi.CreateServiceBindingResponse, error) {
	credentials, err := client.GetConnectionStrings(c.helmClient, instanceID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// TODO Add secrets
	bindingResponse := &brokerapi.CreateServiceBindingResponse{
		Credentials: credentials,
	}
	return bindingResponse, nil
}

// UnBind
func (c *helmController) UnBind(instanceID string, bindingID string) error {
	return nil
}
