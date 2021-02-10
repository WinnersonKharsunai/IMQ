package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WinnersonKharsunai/IMQ/imq-server/cmd/handler"
	"github.com/WinnersonKharsunai/IMQ/imq-server/config"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	server "github.com/WinnersonKharsunai/IMQ/imq-server/pkg/imq/conn-manager"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfgs.DbUserName, cfgs.DbPassword, cfgs.DbHost, cfgs.DbPort, cfgs.DbName)
	db, err := storage.NewMysqlDB(dsn, log)
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", cfgs.ImqServerHost, cfgs.ImqServerPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start imq server on %s: %v", addr, err)
	}
	defer lis.Close()

	svc := handler.NewSevice(log, db)

	s := server.NewServer(cfgs.MaxClient, log, svc)

	serverError := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	defer close(shutdown)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	log.Infof("imq-server running on port: %v", addr)
	go func() {
		s.ListenAndServe(lis)
	}()

	// shutdown
	select {
	case err := <-serverError:
		log.Fatalf("main: fatal error", "error", err)
	case <-shutdown:
		log.Infof("main: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(cfgs.ShutdownGrace))
		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			if err := s.Shutdown(ctx); err != nil {
				log.Warnf("graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
