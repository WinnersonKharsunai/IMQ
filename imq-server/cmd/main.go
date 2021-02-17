package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WinnersonKharsunai/IMQ/imq-server/cmd/handler"
	"github.com/WinnersonKharsunai/IMQ/imq-server/cmd/routes"
	"github.com/WinnersonKharsunai/IMQ/imq-server/config"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/WinnersonKharsunai/IMQ/imq-server/pkg/protocol"
	"github.com/WinnersonKharsunai/IMQ/imq-server/pkg/server"
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
		log.Fatalf("main: failed to get configs: %v", err)
	}

	// configure application and start service

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfgs.DbUserName, cfgs.DbPassword, cfgs.DbHost, cfgs.DbPort, cfgs.DbName)
	db, err := storage.NewMysqlDB(dsn, log)
	if err != nil {
		log.Fatal(err)
	}

	svc := handler.NewSevice(log, db)

	p := protocol.NewProtocol()

	r := routes.NewRouter(svc, p)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqServerHost, cfgs.ImqServerPort)
	s := server.NewServer(log, addr, r)

	log.Infof("main: imq-server running on port: %v", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("main: fatal error: %v", err)
	}

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer close(shutdown)

	select {
	case <-shutdown:
		log.Infof("main: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(cfgs.ShutdownGrace))
		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			if err := s.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
