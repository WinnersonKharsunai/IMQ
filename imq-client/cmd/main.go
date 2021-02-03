package main

import (
	"fmt"
	"net"

	"github.com/WinnersonKharsunai/IMQ/imq-client/cmd/services"
	"github.com/WinnersonKharsunai/IMQ/imq-client/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// load configuration from environment variables
	var cfgs config.Settings
	err := envconfig.Process("", &cfgs)
	if err != nil {
		log.Fatalf("failed to get configs: %v", err)
	}

	// configure application and start service

	imqService := services.NewImqClientService(log)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)
	con, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start connection: %v", err)
	}
	defer con.Close()

	imqService.HandleImqRequest(con)
}
