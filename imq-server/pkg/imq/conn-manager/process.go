package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	sendMessage = "SendMessage"
)

func (s *Server) processWorker() {
	shutdown := false
	for !shutdown {
		select {
		case <-s.processWorkerShutdown:
			shutdown = true
		default:
			con, ok := <-s.newConnection
			for ok {
				data, err := bufio.NewReader(con).ReadBytes('\n')
				if err != nil {
					s.log.Errorf("failed to read client request: %v", err)
				}
				var (
					request  Request
					response string
				)

				if strings.Contains(string(data), "{") {

					if err := json.Unmarshal(data, &request); err != nil {
						s.log.Error(err)
					}

					if err = s.validate.Header(request.Header.Version, request.Header.Method); err != nil {
						s.log.Error(err)
						continue
					}

					resp, err := s.handleRequest(context.Background(), request.Header.Method, request.Body)
					if err != nil {
						s.log.Error(err)
					}
					rb, err := json.Marshal(resp)
					if err != nil {
						s.log.Error(err)
					}
					response = string(rb) + "\n"
				}

				_, err = fmt.Fprintf(con, "%v", response)
				if err != nil {
					s.log.Error(err)
				}
			}
			con.Close()
		}
	}
	s.processWorkerWg.Done()
}

func (s *Server) handleRequest(ctx context.Context, method string, body interface{}) (interface{}, error) {

	switch method {
	case sendMessage:
		sendMessageRequest := &SendMessageRequest{}
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, &sendMessageRequest); err != nil {
			return nil, err
		}

		sendMessageResponse, err := s.imqSvc.SendMessage(ctx, sendMessageRequest)
		if err != nil {
			return nil, err
		}

		return sendMessageResponse, nil
	}
	return nil, nil
}
