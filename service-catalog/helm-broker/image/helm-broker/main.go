package main

import (
	"flag"
	"helm-broker/controller"
	"helm-broker/server"
	"log"
)

var (
	configFilePath string
	port           int
)

func init() {
	flag.StringVar(&configFilePath, "config", "/config/broker-config.json",
		"Location of the config file")
	flag.IntVar(&port, "port", 8080, "Listen port")
}

func main() {
	log.Println("Starting Helm broker...")
	flag.Parse()
	config, err := controller.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	helmController, err := controller.CreateController(config)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(8080, helmController)
}
