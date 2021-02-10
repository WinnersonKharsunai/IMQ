package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/WinnersonKharsunai/IMQ/imq-client/internal/domain"
	"github.com/sirupsen/logrus"
)

// ImqClientService is the reciever type
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

	clientName, err := getClientName()
	if err != nil {
		imq.log.Errorf("failed to get client name: %v", err)
	}

	for {
		msg, err := getMessage()
		if err != nil {
			imq.log.Errorf("failed to read input message: %v", err)
		}

		request := domain.PrepareClientRequest(clientName, msg)

		requestBody, err := marshalRequest(request)
		if err != nil {
			imq.log.Errorf("failed to marshal request: %v", err)
		}

		fmt.Println("request sent:", requestBody)
		_, err = fmt.Fprintf(con, requestBody)
		if err != nil {
			imq.log.Errorf("failed to send request to server: %v", err)
		}

		response, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			imq.log.Errorf("failed to recieved respose: %v", err)
		}

		imq.log.Infof("recieved: %v", response)
	}
}

func getClientName() (string, error) {

	fmt.Print("\nEnter your name: ")
	reader := bufio.NewReader(os.Stdin)

	clientName, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return clientName, nil
}

func getMessage() (string, error) {

	fmt.Print("\nEnter message: ")
	reader := bufio.NewReader(os.Stdin)

	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return msg, nil
}

func marshalRequest(request domain.Request) (string, error) {

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	requestbody := string(requestBytes) + "\n"

	return requestbody, nil
}
