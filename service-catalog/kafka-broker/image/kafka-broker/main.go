package main

import (
	"code.cloudfoundry.org/lager"
	"flag"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/server"
	"kafka-broker/controller"
	"log"
)

var (
	configFilePath string
	port           int
	logLevels      = map[string]lager.LogLevel{
		"DEBUG": lager.DEBUG,
		"INFO":  lager.INFO,
		"ERROR": lager.ERROR,
		"FATAL": lager.FATAL,
	}
)

func init() {
	flag.StringVar(&configFilePath, "config", "/config/broker-config.json",
		"Location of the config file")
	flag.IntVar(&port, "port", 8080, "Listen port")
}

func main() {
	log.Println("Starting Kaka broker...")
	flag.Parse()
	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	controller, err := controller.CreateController(config.KafkaConfig)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(port, controller)
}
