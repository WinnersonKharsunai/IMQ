package server

import (
	"context"
	"net"
	"sync"

	protocol "github.com/WinnersonKharsunai/IMQ/imq-server/pkg/imq/conn-protocol"
	"github.com/sirupsen/logrus"
)

// Server is the reciever type
type Server struct {
	maxClient             int
	imqSvc                ImqServerServiceIF
	validate              protocol.Validator
	log                   *logrus.Logger
	newConnection         chan net.Conn
	processedMessages     chan interface{}
	processWorkerShutdown chan struct{}
	processWorkerWg       sync.WaitGroup
	listenerWorkerWg      sync.WaitGroup
}

// NewServer is the factory function for Server
func NewServer(maxClient int, log *logrus.Logger, svc ImqServerServiceIF) *Server {
	return &Server{
		maxClient:             maxClient,
		log:                   log,
		imqSvc:                svc,
		newConnection:         make(chan net.Conn),
		processedMessages:     make(chan interface{}),
		processWorkerShutdown: make(chan struct{}),
	}
}

// ListenAndServe initaites all the background workers and start serving
func (s *Server) ListenAndServe(lis net.Listener) {
	s.processWorkerWg.Add(s.maxClient)

	for i := 0; i < s.maxClient; i++ {
		go s.processWorker()
	}

	s.listenerWorkerWg.Add(1)
	go s.listener(lis)
}

// Shutdown stops all the background workers
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(s.processWorkerShutdown)
		s.listenerWorkerWg.Wait()

		close(s.newConnection)
		s.processWorkerWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
