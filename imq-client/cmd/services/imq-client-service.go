package services

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/WinnersonKharsunai/IMQ/imq-client/internal/domain"
	"github.com/sirupsen/logrus"
)

// ImqClientService ...
type ImqClientService struct {
	log *logrus.Logger
}

// NewImqClientService is the factory function for ImqClientService
func NewImqClientService(l *logrus.Logger) *ImqClientService {
	return &ImqClientService{
		log: l,
	}
}

// HandleImqRequest ...
func (imq *ImqClientService) HandleImqRequest(con net.Conn) {

	clientName := getClientName()

	for {
		msg := getMessage()
		request := domain.PrepareClientRequest(clientName, msg)

		requestBody, err := marshalRequest(request)
		if err != nil {
			imq.log.Errorf("failed to marshal request: %v", err)
		}

		_, err = fmt.Fprintf(con, requestBody+"\n")
		if err != nil {
			imq.log.Errorf("failed to send request to server: %v", err)
		}
	}
}

func getClientName() string {

	var name string

	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)

	return name
}

func getMessage() string {

	var msg string

	fmt.Print("Enter message: ")
	fmt.Scanln(&msg)

	return msg
}

func marshalRequest(request domain.Request) (string, error) {

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	requestbody := string(requestBytes)

	return requestbody, nil
}
