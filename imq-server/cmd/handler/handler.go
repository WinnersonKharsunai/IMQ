package handler

import (
	"context"
	"fmt"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/domain"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	server "github.com/WinnersonKharsunai/IMQ/imq-server/pkg/imq/conn-manager"
	"github.com/sirupsen/logrus"
)

// ServerService is the reciever type
type ServerService struct {
	db  storage.DatabaseIF
	log *logrus.Logger
}

// NewSevice is the factory function for ServerService
func NewSevice(log *logrus.Logger, db storage.DatabaseIF) server.ImqServerServiceIF {
	return &ServerService{
		db:  db,
		log: log,
	}
}

// SendMessage recieve and persists client message
func (s *ServerService) SendMessage(ctx context.Context, in *server.SendMessageRequest) (*server.SendMessageResponse, error) {

	fmt.Println("request recieved:", in)

	request := &domain.Request{
		ClientName: in.ClientName,
		Message: domain.Message{
			Data:      in.Data,
			CreatedAt: in.CreatedAt,
			ExpireAt:  in.ExpireAt,
		},
	}

	if err := domain.ProcessSendMessage(ctx, s.log, s.db, request); err != nil {
		s.log.Errorf("failed to process client request: %v", err)
	}

	return &server.SendMessageResponse{}, nil
}
