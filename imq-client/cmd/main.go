package main

import (
	"fmt"
	"net"

	"github.com/WinnersonKharsunai/IMQ/imq-client/cmd/handler"
	"github.com/WinnersonKharsunai/IMQ/imq-client/config"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

func main() {
	// initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// load configuration from environment variables
	cfgs := config.Settings{}
	if err := env.Parse(&cfgs); err != nil {
		log.Fatalf("failed to get configs: %v", err)
	}

	// configure application and start service

	imqService := handler.NewImqClientService(log)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)
	con, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start connection: %v", err)
	}
	defer con.Close()

	for {
		imqService.HandleImqRequest(con)
	}
}
