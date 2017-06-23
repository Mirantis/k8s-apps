package main

import (
	"flag"
	"fmt"
	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1alpha1"
	"io/ioutil"
	"path/filepath"
	"service-catalog-client/client"
	"time"
)

const (
	BrokerName      = "helm-broker"
	BrokerURL       = "http://helm-broker-helm-broker.helm-broker.svc.cluster.local"
	ChartRepo       = "mirantisworkloads.storage.googleapis.com"
	InstanceName    = "helm-instance"
	Chart           = "zookeeper"
	InstanceVersion = "1.1.0"
	Namespace       = "test-ns"
	BindingName     = "helm-binding"
	ValuesPath      = "./values.json"
)

var (
	brokerName      string
	instanceName    string
	chart           string
	instanceVersion string
	namespace       string
	bindingName     string
)

func init() {
	flag.StringVar(&brokerName, "broker", BrokerName, "Broker name")
	flag.StringVar(&instanceName, "instance", InstanceName, "Instance name")
	flag.StringVar(&chart, "chart", Chart, "Chart name")
	flag.StringVar(&instanceVersion, "version", InstanceVersion, "Chart Version")
	flag.StringVar(&namespace, "namespace", Namespace, "")
	flag.StringVar(&bindingName, "binding", BindingName, "Binding name")
}

func main() {
	// Create client with current k8s config and context
	cli, _ := client.NewClient()
	CreateBroker(*cli)
	PrintServiceClasses(*cli)
	CreateInstance(*cli)
	Bind(*cli)
}

func CreateBroker(cli client.CatalogClient) {
	br, err := cli.CreateBroker(brokerName, BrokerURL)
	fmt.Printf("Broker %s creation is in progress...\n", br.Name)
	if err != nil {
		panic(err.Error())
	}
	state := v1alpha1.ConditionStatus("Unknown")
	for {
		br, err := cli.GetBroker(brokerName)
		if err != nil {
			panic(err.Error())
		}
		conditions := br.Status.Conditions
		if len(conditions) != 0 {
			state = br.Status.Conditions[0].Status
			break
		}
	}
	if state != "True" {
		fmt.Printf("Broker %s hasn't been created\n", br.Name)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Printf("Broker %s has been created\n", br.Name)
	}
}

func PrintServiceClasses(cli client.CatalogClient) {
	sclasses, err := cli.GetServiceClasses()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("The following Service Classes were created:")
	for _, sc := range sclasses.Items {
		fmt.Println(sc.Name)
	}
}

func CreateInstance(cli client.CatalogClient) {
	path, _ := filepath.Abs(ValuesPath)
	values, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	instance, err := cli.CreateInstance(instanceName, chart, instanceVersion, ChartRepo, namespace, values)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(instance.Name)

	state := v1alpha1.ConditionStatus("Unknown")
	for i := 0; i < 100; i++ {
		in, err := cli.GetInstance(instanceName, namespace)
		if err != nil {
			panic(err.Error())
		}
		conditions := in.Status.Conditions
		if len(conditions) == 0 {
			time.Sleep(time.Second * 3)
			continue
		}
		state = in.Status.Conditions[0].Status
		if state == "True" {
			break
		}
		time.Sleep(time.Second * 3)
	}
	if state != "True" {
		fmt.Printf("Instance %s hasn't been created\n", instance.Name)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Printf("Instance %s has been created\n", instance.Name)
	}
}

func Bind(cli client.CatalogClient) {
	binding, err := cli.Bind(bindingName, instanceName, namespace)
	if err != nil {
		panic(err.Error())
	}
	state := v1alpha1.ConditionStatus("Unknown")
	for {
		br, err := cli.GetBinding(bindingName, namespace)
		if err != nil {
			panic(err.Error())
		}
		conditions := br.Status.Conditions
		if len(conditions) != 0 {
			state = br.Status.Conditions[0].Status
			break
		}
	}
	if state == "True" {
		fmt.Printf("Binding %s has been created\n", binding.Name)
	} else {
		fmt.Printf("Binding %s hasn't been created\n", binding.Name)
	}
}
