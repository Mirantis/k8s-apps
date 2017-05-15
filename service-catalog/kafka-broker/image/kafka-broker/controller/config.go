package controller

import (
	"errors"
)

type Config struct {
	Topic             string   `json:"topic"`
	KafkaBrokers      []string `json:"brokers"`
	ZookeeperServer   string   `json:"zookeeperServer"`
	Partitions        int      `json:"partitions"`
	ReplicationFactor int      `json:"replicationFactor"`
}

func (c Config) Validate() error {
	if len(c.KafkaBrokers) == 0 {
		return errors.New("Must provide a non-empty Broker list")
	}
	if c.Topic == "" {
		return errors.New("Must provide a non-empty Topic")
	}
	if c.ZookeeperServer == "" {
		return errors.New("Must provide a non-empty ZookeeperServer")
	}
	if c.Partitions <= 0 {
		return errors.New("Must provide a non-zero Partitions")
	}
	if c.ReplicationFactor <= 0 {
		return errors.New("Must provide a non-zero ReplicationFactor")
	}
	return nil
}
