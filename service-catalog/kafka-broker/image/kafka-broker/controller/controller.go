package controller

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"log"
)

type kafkaController struct {
	kafkaClient       sarama.Client
	topic             string
	kafkaBrokers      []string
	topicManager      TopicManager
	partitions        int
	replicationFactor int
	zookeeperServer   string
}

// CreateController returns a Kafka Broker Controller
func CreateController(c Config) (controller.Controller, error) {

	kafkaConfig := sarama.NewConfig()
	kafkaClient, err := sarama.NewClient(c.KafkaBrokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	topicManager, err := CreateTopicManager(c.ZookeeperServer, kafkaClient)
	if err != nil {
		return nil, err
	}

	return &kafkaController{
		kafkaClient:       kafkaClient,
		topic:             c.Topic,
		kafkaBrokers:      c.KafkaBrokers,
		partitions:        c.Partitions,
		replicationFactor: c.ReplicationFactor,
		topicManager:      topicManager,
		zookeeperServer:   c.ZookeeperServer,
	}, nil
}

// Catalog returns the Kafka Broker catalog entries
func (c *kafkaController) Catalog() (*brokerapi.Catalog, error) {
	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				Name:        "kafka",
				ID:          "effc24e9-0a40-4cd3-a3e4-0435ae156a43",
				Description: "Kafka",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "f6d1d1ba-ab4a-4d68-973b-6e5f6c8723ca",
						Description: "Kafka",
						Free:        true,
					},
				},
				Bindable: true,
			},
		},
	}, nil
}

// CreateServiceInstance
func (c *kafkaController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (
	*brokerapi.CreateServiceInstanceResponse, error) {
	topicName := c.topicName(id)
	topic, err := c.topicManager.CreateTopic(topicName, c.partitions, c.replicationFactor)
	if err != nil {
		log.Println("Topic creation was failed with the following error: ")
		log.Println(err)
		return nil, err
	}
	log.Println("Topic has been successfully created: " + topic)
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

// GetServiceInstance
func (c *kafkaController) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

// RemoveServiceInstance
func (c *kafkaController) RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error) {
	topicName := c.topicName(id)
	topic, err := c.topicManager.DeleteTopic(topicName)
	if err != nil {
		log.Println("Topic deleting was failed with the following error: ")
		log.Println(err)
		return nil, err
	}
	log.Println("Topic has been successfully removed: " + topic)
	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

// Bind
func (c *kafkaController) Bind(instanceID string, bindingID string, req *brokerapi.BindingRequest) (
	*brokerapi.CreateServiceBindingResponse, error) {
	topicName := c.topicName(instanceID)
	bindingResponse := &brokerapi.CreateServiceBindingResponse{
		Credentials: brokerapi.Credential{
			"name":            topicName,
			"brokers":         c.kafkaBrokers,
			"zookeeperServer": c.zookeeperServer,
		},
	}
	return bindingResponse, nil
}

// UnBind
func (c *kafkaController) UnBind(instanceID string, bindingID string) error {
	return nil
}

func (c *kafkaController) topicName(instanceID string) string {
	return fmt.Sprintf("%s-%s", c.topic, instanceID)
}
