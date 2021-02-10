package storage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WinnersonKharsunai/IMQ/imq-server/internal/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var clientData *storage.ClientData

func init() {
	now := time.Now().UTC()
	clientData = &storage.ClientData{
		Data:      "test data",
		CreatedAt: now.Format("2006-01-02 15:04:05"),
		ExpiresAt: now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
	}
}

func TestStoreClientData_Failed(t *testing.T) {

	clientID := "testId123"
	expectedErr := errors.New("failed to insert clientData")
	mock, db := mysqlMock()

	stmt := `INSERT INTO`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.StoreClientData(context.Background(), clientID, clientData)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestStoreClientData_Success(t *testing.T) {

	clientID := "testId123"
	mock, db := mysqlMock()

	stmt := `INSERT INTO`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.StoreClientData(context.Background(), clientID, clientData)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestPrepareTableForNewClient_Failed(t *testing.T) {

	clientID := uuid.New().String()
	expectedErr := errors.New("failed to create new table")
	mock, db := mysqlMock()

	stmt := `CREATE TABLE`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.PrepareTableForNewClient(context.Background(), clientID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestPrepareTableForNewClient_Success(t *testing.T) {

	clientID := uuid.New().String()
	mock, db := mysqlMock()

	stmt := `CREATE TABLE`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.PrepareTableForNewClient(context.Background(), clientID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestStoreNewClientID_Failed(t *testing.T) {

	clientName := "John"
	clientID := uuid.New().String()
	expectedErr := errors.New("failed to insert")
	mock, db := mysqlMock()

	stmt := `INSERT INTO Clients \(clientName,clientId\) VALUES\(\?,\?\)`
	mock.ExpectExec(stmt).WithArgs(clientName, clientID).WillReturnError(expectedErr)

	err := db.StoreNewClientID(context.Background(), clientName, clientID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestStoreNewClientID_Success(t *testing.T) {

	clientName := "John"
	clientID := uuid.New().String()
	mock, db := mysqlMock()

	stmt := `INSERT INTO Clients \(clientName,clientId\) VALUES\(\?,\?\)`
	mock.ExpectExec(stmt).WithArgs(clientName, clientID).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.StoreNewClientID(context.Background(), clientName, clientID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestGetClientID_Failed(t *testing.T) {

	clientName := "John"
	expectedErr := errors.New("failed to get clientId")
	mock, db := mysqlMock()

	stmt := `SELECT clientId FROM Clients where clientName = \?`
	mock.ExpectQuery(stmt).WithArgs(clientName).WillReturnError(expectedErr)

	clientId, err := db.GetClientID(context.Background(), clientName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
	if clientId != "" {
		t.Fatalf("clientId should be empty")
	}
}

func TestGetClientID_Success(t *testing.T) {

	clientName := "John"
	mock, db := mysqlMock()

	columns := []string{"client_name"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("John")

	stmt := `SELECT clientId FROM Clients where clientName = \?`
	mock.ExpectQuery(stmt).WithArgs(clientName).WillReturnRows(rows)

	clientId, err := db.GetClientID(context.Background(), clientName)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
	if clientId == "" {
		t.Fatalf("clientId should not be empty")
	}
}

func mysqlMock() (sqlmock.Sqlmock, storage.MysqlDB) {

	dbCxn, mock, _ := sqlmock.New()
	db := storage.MysqlDB{
		Cxn: dbCxn,
		Log: logrus.New(),
	}

	return mock, db
}
