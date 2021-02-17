package server

import (
	"context"
	"fmt"
	"time"
)

func (s *Server) listenerWorker() {
	shutdown := false
	for !shutdown {
		select {
		case <-s.shutdown:
			shutdown = true
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			go s.accept(ctx)

			select {
			case <-ctx.Done():
			case ch := <-s.listenerCh:
				if ch.err != nil {
					fmt.Println("got connection error:", ch.err)
				} else {
					go s.processWorker(ch.con)
				}
			}
		}
	}
	s.listenerWg.Done()
}

func (s *Server) accept(ctx context.Context) {
	con, err := s.listener.Accept()
	s.listenerCh <- connection{
		con: con,
		err: err,
	}
}
