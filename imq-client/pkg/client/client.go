package client

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/WinnersonKharsunai/IMQ/imq-client/cmd/handler"
)

// Client is the receiver type
type Client struct {
	addr      string
	svc       *handler.ClientService
	processCh chan net.Conn
	processWg sync.WaitGroup
}

// NewClient is the factory function for client
func NewClient(addr string, svc *handler.ClientService) *Client {
	return &Client{
		addr:      addr,
		svc:       svc,
		processCh: make(chan net.Conn),
	}
}

// Dial connects to server
func (c *Client) Dial() error {

	con, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Fatalf("failed to start connection: %v", err)
		return err
	}

	c.processWg.Add(1)
	go c.processWorker()

	c.processCh <- con

	return nil
}

func (c *Client) processWorker() {

	for con, chanOpen := <-c.processCh; chanOpen; {
		defer con.Close()
		c.svc.HandleImqRequest(con)
	}
	c.processWg.Done()
}

// Shutdown gracefully shutdown the client
func (c *Client) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(c.processCh)
		c.processWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
