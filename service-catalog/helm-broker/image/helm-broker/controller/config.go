package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ChartUrls  []string `json:"chartUrls"`
	TillerHost string   `json:"tillerHost"`
}

func LoadConfig(configFile string) (config Config, err error) {
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
	if len(c.ChartUrls) == 0 {
		return errors.New("Must provide a non-empty Chart URL list")
	}
	if c.TillerHost == "" {
		return errors.New("Must provide a non-empty Tiller Host")
	}
	return nil
}
