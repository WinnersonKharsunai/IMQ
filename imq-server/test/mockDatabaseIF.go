package test

import (
	"context"

	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
)

// MockDatabaseIF is a struct for mocking DatabaseIF
type MockDatabaseIF struct {
	Mock
	storage.DatabaseIF
}

// GetClientID mocks on GetClientID.GetClientID
func (m *MockDatabaseIF) GetClientID(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

// StoreNewClientID mocks on GetClientID.StoreNewClientID
func (m *MockDatabaseIF) StoreNewClientID(ctx context.Context, clientName, clientID string) error {
	args := m.Called(ctx, clientName, clientID)
	return args.Error(0)
}

// PrepareTableForNewClient mocks on GetClientID.PrepareTableForNewClient
func (m *MockDatabaseIF) PrepareTableForNewClient(ctx context.Context, clientID string) error {
	args := m.Called(ctx, clientID)
	return args.Error(0)
}

// StoreClientData mocks on GetClientID.StoreClientData
func (m *MockDatabaseIF) StoreClientData(ctx context.Context, clientID string, cd *storage.ClientData) error {
	args := m.Called(ctx, clientID, cd)
	return args.Error(0)
}
