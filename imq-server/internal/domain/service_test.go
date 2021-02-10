package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/domain"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/WinnersonKharsunai/IMQ/imq-server/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var request *domain.Request

func init() {
	now := time.Now().UTC()
	request = &domain.Request{
		ClientName: "John",
		Message: domain.Message{
			Data:      "test data",
			CreatedAt: now.Format("2006-01-02 15:04:05"),
			ExpireAt:  now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
		},
	}
}

func TestProcessImqRequest_GetClientID_Failed(t *testing.T) {

	expectedErr := errors.New("failed to get clientId")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetClientID).When(mock.Anything, request.ClientName).Return("", expectedErr)

	err := domain.ProcessSendMessage(context.Background(), &logrus.Logger{}, mockDb, request)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestProcessImqRequest_StoreNewClientID_Failed(t *testing.T) {

	expectedErr := errors.New("failed to store new client")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetClientID).When(mock.Anything, request.ClientName).Return("", nil)
	mockDb.Given(storage.DatabaseIF.StoreNewClientID).When(mock.Anything, request.ClientName, mock.Anything).Return(expectedErr)

	err := domain.ProcessSendMessage(context.Background(), &logrus.Logger{}, mockDb, request)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestProcessImqRequest_PrepareTableForNewClient_Failed(t *testing.T) {

	expectedErr := errors.New("failed to prepare new table for client")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetClientID).When(mock.Anything, request.ClientName).Return("", nil)
	mockDb.Given(storage.DatabaseIF.StoreNewClientID).When(mock.Anything, request.ClientName, mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.PrepareTableForNewClient).When(mock.Anything, mock.Anything).Return(expectedErr)

	err := domain.ProcessSendMessage(context.Background(), &logrus.Logger{}, mockDb, request)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestProcessImqRequest_StoreClientData_Failed(t *testing.T) {

	expectedErr := errors.New("failed to store client data")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetClientID).When(mock.Anything, request.ClientName).Return("", nil)
	mockDb.Given(storage.DatabaseIF.StoreNewClientID).When(mock.Anything, request.ClientName, mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.PrepareTableForNewClient).When(mock.Anything, mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.StoreClientData).When(mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	err := domain.ProcessSendMessage(context.Background(), &logrus.Logger{}, mockDb, request)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestProcessImqRequest_Success(t *testing.T) {

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetClientID).When(mock.Anything, request.ClientName).Return("", nil)
	mockDb.Given(storage.DatabaseIF.StoreNewClientID).When(mock.Anything, request.ClientName, mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.PrepareTableForNewClient).When(mock.Anything, mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.StoreClientData).When(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := domain.ProcessSendMessage(context.Background(), &logrus.Logger{}, mockDb, request)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}
