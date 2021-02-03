package services

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/domain"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/sirupsen/logrus"
)

// ImqServerService ...
type ImqServerService struct {
	log *logrus.Logger
	db  storage.DatabaseIF
}

// NewImqServerService is the factory function for ImqServerService
func NewImqServerService(db storage.DatabaseIF, l *logrus.Logger) *ImqServerService {
	return &ImqServerService{
		db:  db,
		log: l,
	}
}

// HandleImqRequest ...
func (imq *ImqServerService) HandleImqRequest(ctx context.Context, con net.Conn) {

	defer con.Close()

	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			imq.log.Errorf("failed to read client request: %v", err)
		}

		request, err := unmarshalRequest(data)
		if err != nil {
			imq.log.Errorf("failed to unmarshal client request: %v", err)
		}

		if err := domain.ProcessImqRequest(ctx, imq.log, imq.db, request); err != nil {
			imq.log.Errorf("failed to process client request: %v", err)
		}

		response, err := marshalResponse(*request)
		if err != nil {
			imq.log.Errorf("failed to marshal response: %v", err)
		}

		_, err = con.Write(response)
		if err != nil {
			imq.log.Errorf("failed to send response: %v", err)
		}
	}
}

func unmarshalRequest(requestBody string) (*domain.Request, error) {

	var request domain.Request

	if err := json.Unmarshal([]byte(requestBody), &request); err != nil {
		return nil, err
	}

	return &request, nil
}

func marshalResponse(request domain.Request) ([]byte, error) {

	response := domain.Response{
		Name:      request.ClientName,
		Timestamp: request.Message.CreatedAt,
		Data:      request.Message.Data,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
