package controller

import (
	"github.com/Shopify/sarama"
	"github.com/samuel/go-zookeeper/zk"
	"kafka-broker/utils"
	"time"
)

type TopicManager interface {
	CreateTopic(topic string, partitions int, replicationFactor int) (string, error)
	DeleteTopic(topic string) (string, error)
}

type topicManager struct {
	zookeeperConnection zk.Conn
	kafkaClient         sarama.Client
}

func CreateTopicManager(zkServer string, kafkaClient sarama.Client) (TopicManager, error) {
	zkConnection, _, err := zk.Connect([]string{zkServer}, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	return &topicManager{
		zookeeperConnection: *zkConnection,
		kafkaClient:         kafkaClient,
	}, nil
}

func (tm *topicManager) CreateTopic(topic string, partitions int, replicationFactor int) (string, error) {
	brokerList := tm.kafkaClient.Brokers()
	brokerIDList := make([]int32, len(brokerList))
	for i, broker := range brokerList {
		brokerIDList[i] = broker.ID()
	}
	replicaAssignment, err := utils.AssignReplicasToBrokers(brokerIDList, partitions, replicationFactor)
	if err != nil {
		return topic, err
	}
	err = utils.WriteTopicPartitionAssignment(tm.zookeeperConnection, topic, replicaAssignment)
	return topic, err
}

func (tm *topicManager) DeleteTopic(topic string) (string, error) {
	err := utils.DeleteTopicPartitionAssignment(tm.zookeeperConnection, topic)
	return topic, err
}
