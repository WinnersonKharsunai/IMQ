package handler

import (
	"context"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/domain"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/WinnersonKharsunai/IMQ/imq-server/pkg/services"
	"github.com/sirupsen/logrus"
)

// ServerService is the receiver type
type ServerService struct {
	log *logrus.Logger
	db  storage.DatabaseIF
}

// NewSevice is the factory function for ServerService
func NewSevice(log *logrus.Logger, db storage.DatabaseIF) services.ImqServiceIF {
	return &ServerService{
		log: log,
		db:  db,
	}
}

// SendMessage recieve and persists client message
func (s *ServerService) SendMessage(ctx context.Context, in *services.SendMessageRequest) (*services.SendMessageResponse, error) {

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
		return nil, err
	}

	return &services.SendMessageResponse{
		Data:      in.Data,
		CreatedAt: in.CreatedAt,
		ExpireAt:  in.CreatedAt,
	}, nil
}
