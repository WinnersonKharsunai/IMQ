package server

import (
	"fmt"
	"net"
)

func (s *Server) listener(lis net.Listener) {
	shutdown := false
	for !shutdown {
		select {
		case <-s.processWorkerShutdown:
			fmt.Println("listener closing")
			lis.Close()
			shutdown = true
		default:
			fmt.Println("listener starting")
			conn, err := lis.Accept()
			if err != nil {
				s.log.Warnf("error while accepting client connection: %v", err)
				continue
			}
			fmt.Println("got a new connection")
			s.newConnection <- conn
		}
	}
	s.listenerWorkerWg.Done()
}
