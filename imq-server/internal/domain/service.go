package domain

import (
	"context"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ProcessSendMessage processes the requests sent by client
func ProcessSendMessage(ctx context.Context, log *logrus.Logger, db storage.DatabaseIF, request *Request) error {

	clientID, err := db.GetClientID(ctx, request.ClientName)
	if err != nil {
		return err
	}

	if clientID == "" {
		clientID = uuid.New().String()

		err = db.StoreNewClientID(ctx, request.ClientName, clientID)
		if err != nil {
			return err
		}

		if err := db.PrepareTableForNewClient(ctx, clientID); err != nil {
			return err
		}
	}

	cd := storage.ClientData{
		CreatedAt: request.Message.CreatedAt,
		ExpiresAt: request.Message.ExpireAt,
		Data:      request.Message.Data,
	}

	if err := db.StoreClientData(ctx, clientID, &cd); err != nil {
		return err
	}

	return nil
}
