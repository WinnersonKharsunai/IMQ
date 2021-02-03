package main

import (
	"context"
	"fmt"
	"net"

	"github.com/WinnersonKharsunai/IMQ/imq-server/cmd/services"
	"github.com/WinnersonKharsunai/IMQ/imq-server/config"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfgs.DbUserName, cfgs.DbPassword, cfgs.DbHost, cfgs.DbPort, cfgs.DbName)
	db, err := storage.NewMysqlDB(dsn, log)
	if err != nil {
		log.Fatal(err)
	}

	imqService := services.NewImqServerService(db, log)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqServerHost, cfgs.ImqServerPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start imq server on %s: %v", addr, err)
	}
	defer lis.Close()

	log.Infof("imq server running on port: %v", addr)
	for {
		con, err := lis.Accept()
		if err != nil {
			log.Errorf("error while accepting connection: %v", err.Error())
		}

		go imqService.HandleImqRequest(context.Background(), con)
	}
}
