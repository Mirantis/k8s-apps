package utils

import (
	"encoding/json"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
)

var BrokerTopicsPath = "/brokers/topics"
var DeleteTopicsPath = "/admin/delete_topics"

type zkPath struct {
	Version    int             `json:"version"`
	Partitions map[int][]int32 `json:"partitions"`
}

func WriteTopicConfig() {
	// TODO implement it
}
func WriteTopicPartitionAssignment(zookeeperConnection zk.Conn, topic string,
	partitionReplicaAssignment map[int][]int32) error {
	zkPath := BrokerTopicsPath + "/" + topic
	jsonPartitionData, err := replicaAssignmentZkData(partitionReplicaAssignment)
	if err != nil {
		return err
	}
	err = createPersistentPath(zookeeperConnection, zkPath, jsonPartitionData)
	return err
}

func DeleteTopicPartitionAssignment(zookeeperConnection zk.Conn, topic string) error {
	zkPath := DeleteTopicsPath + "/" + topic
	err := createPersistentPath(zookeeperConnection, zkPath, []byte{})
	return err
}

func AssignReplicasToBrokers(brokerList []int32, partitions int, replicationFactor int) (map[int][]int32, error) {

	if partitions <= 0 {
		return nil, errors.New("Partitions count must be larger than 0")
	}
	if replicationFactor <= 0 {
		return nil, errors.New("Replication factor must be larger than 0")
	}
	if replicationFactor > int(len(brokerList)) {
		return nil, errors.New("Replication factor: " + string(replicationFactor) +
			" larger than available brokers: " + string(len(brokerList)))
	}
	result := make(map[int][]int32)

	currentPartitionId := 0
	nextReplicaShift := 0
	firstReplicaIndex := 0

	for i := 0; i < partitions; i++ {
		replicaList := make([]int32, replicationFactor)
		if currentPartitionId > 0 && (currentPartitionId%len(brokerList) == 0) {
			nextReplicaShift += 1
		}
		replicaList[0] = brokerList[firstReplicaIndex]
		for j := 0; j < replicationFactor-1; j++ {
			replicaList[j+1] = brokerList[replicaIndex(
				firstReplicaIndex, nextReplicaShift, j, len(brokerList))]
		}
		result[currentPartitionId] = replicaList
		currentPartitionId += 1
		if firstReplicaIndex != len(brokerList)-1 {
			firstReplicaIndex++
		} else {
			firstReplicaIndex = 0
		}

	}
	return result, nil
}

func replicaAssignmentZkData(partitionReplicaAssignment map[int][]int32) ([]byte, error) {
	zkPath := &zkPath{
		Version:    1,
		Partitions: partitionReplicaAssignment,
	}
	result, err := json.Marshal(zkPath)
	return result, err
}

func replicaIndex(firstReplicaIndex int, secondReplicaShift int, replicaIndex int, nBrokers int) int {
	shift := 1 + (secondReplicaShift+replicaIndex)%(nBrokers-1)
	return (firstReplicaIndex + shift) % nBrokers
}

func createPersistentPath(zookeeperConnection zk.Conn, zkPath string, jsonPartitionData []byte) error {
	_, err := zookeeperConnection.Create(zkPath, jsonPartitionData, 0, zk.WorldACL(zk.PermAll))
	return err
}
