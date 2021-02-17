package server

import (
	"context"
	"net"
	"sync"

	"github.com/WinnersonKharsunai/IMQ/imq-server/cmd/routes"
	"github.com/sirupsen/logrus"
)

// Server is the receiver type
type Server struct {
	address    string
	log        *logrus.Logger
	routes     *routes.Router
	listener   net.Listener
	listenerCh chan connection
	listenerWg sync.WaitGroup
	shutdown   chan struct{}
}

// NewServer is the factory function for Server type
func NewServer(log *logrus.Logger, addr string, routes *routes.Router) *Server {
	return &Server{
		log:        log,
		address:    addr,
		routes:     routes,
		listenerCh: make(chan connection),
		shutdown:   make(chan struct{}),
	}
}

// ListenAndServe start the server and serve multiple clients
func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.listener = lis

	s.listenerWg.Add(1)
	go s.listenerWorker()

	return nil
}

// Shutdown gracefully shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(s.shutdown)
		s.listenerWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
