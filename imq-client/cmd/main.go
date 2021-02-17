package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WinnersonKharsunai/IMQ/imq-client/cmd/handler"
	"github.com/WinnersonKharsunai/IMQ/imq-client/config"
	"github.com/WinnersonKharsunai/IMQ/imq-client/pkg/client"
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

	svc := handler.NewClientService(log)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)

	//log.Infof("main: imq-client dialing on port: %s", addr)
	c := client.NewClient(addr, svc)

	if err := c.Dial(); err != nil {
		log.Fatalf("failed to start connection: %v", err)
	}

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer close(shutdown)

	select {
	case <-shutdown:
		log.Infof("\nmain: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(cfgs.ShutdownGrace))
		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			if err := c.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
