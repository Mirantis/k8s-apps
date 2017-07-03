package client

import (
	"encoding/json"
	"flag"
	v1alpha1 "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1alpha1"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type CatalogClient struct {
	ClientSet *clientset.Clientset
}

func NewClient() (*CatalogClient, error) {

	var kubeConfig *string
	if home := os.Getenv("HOME"); home != "" {
		kubeConfig = flag.String(
			"kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	catalogCient, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &CatalogClient{
		ClientSet: catalogCient,
	}, nil
}

func (cli *CatalogClient) CreateBroker(name, url string) (*v1alpha1.Broker, error) {
	brokerObj := &v1alpha1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.BrokerSpec{
			URL: url,
		},
	}
	broker, err := cli.ClientSet.Brokers().Create(brokerObj)
	return broker, err
}

func (cli *CatalogClient) GetBroker(name string) (*v1alpha1.Broker, error) {
	return cli.ClientSet.Brokers().Get(name, metav1.GetOptions{})
}

func (cli *CatalogClient) GetServiceClasses() (*v1alpha1.ServiceClassList, error) {
	return cli.ClientSet.ServiceClasses().List(metav1.ListOptions{})
}

type Parameters struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
	Repo      string `json:"repo"`
	Values    Values `json:"values"`
}

type Values map[string]interface{}

func (cli *CatalogClient) CreateInstance(name, chart, version, repo, namespace string, values []byte) (*v1alpha1.Instance, error) {
	var valuesMap Values
	if err := json.Unmarshal(values, &valuesMap); err != nil {
		panic(err)
	}

	param := Parameters{
		Name:      chart,
		Namespace: namespace,
		Version:   version,
		Repo:      repo,
		Values:    valuesMap,
	}
	data, err := json.Marshal(param)
	if err != nil {
		panic(err.Error())
	}
	raw := &runtime.RawExtension{
		Raw: data,
	}
	if err != nil {
		panic(err.Error())
	}
	instanceObj := &v1alpha1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.InstanceSpec{
			ServiceClassName: chart + "." + version + "." + repo,
			PlanName:         "default",
			Parameters:       raw,
		},
	}
	return cli.ClientSet.Instances(namespace).Create(instanceObj)
}

func (cli *CatalogClient) GetInstance(name, namespace string) (*v1alpha1.Instance, error) {
	return cli.ClientSet.Instances(namespace).Get(name, metav1.GetOptions{})
}

func (cli *CatalogClient) Bind(name, instanceName, namespace string) (*v1alpha1.Binding, error) {
	bindingObj := &v1alpha1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.BindingSpec{
			InstanceRef: v1.LocalObjectReference{
				Name: instanceName,
			},
			SecretName: instanceName + "-secret",
		},
	}
	binding, err := cli.ClientSet.Bindings(namespace).Create(bindingObj)
	return binding, err
}

func (cli *CatalogClient) GetBinding(name, namespace string) (*v1alpha1.Binding, error) {
	return cli.ClientSet.Bindings(namespace).Get(name, metav1.GetOptions{})
}
