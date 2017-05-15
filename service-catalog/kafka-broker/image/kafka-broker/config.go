package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"kafka-broker/controller"
	"os"
)

type Config struct {
	LogLevel    string            `json:"log_level"`
	KafkaConfig controller.Config `json:"kafka_config"`
}

func LoadConfig(configFile string) (config *Config, err error) {
	if configFile == "" {
		return config, errors.New("Must provide a config file")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	if err = config.Validate(); err != nil {
		return config, fmt.Errorf("Validating config contents: %s", err)
	}

	return config, nil
}

func (c Config) Validate() error {
	if err := c.KafkaConfig.Validate(); err != nil {
		return fmt.Errorf("Validating controller configuration: %s", err)
	}

	return nil
}
