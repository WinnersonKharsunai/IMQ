package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
)

func (s *Server) processWorker(con net.Conn) {

	s.log.Infof("processWorker: start serving new client: %v", con.LocalAddr().String())

	done := false
	for !done {
		request, err := readFromConn(con)
		if err != nil {
			s.log.Errorf("processWorker: connection closed: %v", err)
			done = true
			continue
		}

		resp := s.routes.RequestRouter(context.Background(), request)
		if err != nil {
			s.log.Errorf("processWorker: failed to process request: %v", err)
		}

		if err = writeToConn(con, resp); err != nil {
			s.log.Errorf("processWorker: failed to write response to client: %v", err)
		}

		s.log.Infof("processWorker: completed %v client request with response %+v", con.LocalAddr().String(), resp)
	}
	defer con.Close()
}

func readFromConn(con net.Conn) (interface{}, error) {
	var r interface{}

	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r, nil
}

func writeToConn(con net.Conn, response interface{}) error {

	b, err := json.Marshal(response)
	if err != nil {
		return err
	}

	r := string(b)

	_, err = fmt.Fprintf(con, "%v\n", r)
	if err != nil {
		return err
	}

	return nil
}
