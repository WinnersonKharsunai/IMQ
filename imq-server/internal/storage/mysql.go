package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var mysqlDriver = "mysql"

// DatabaseIF is an interface for messaging queue
type DatabaseIF interface {
	Connect() error
	Test() error
	GetClientID(ctx context.Context, name string) (string, error)
	GenerateClientID(ctx context.Context, name string) (string, error)
	PrepareTableForNewClient(ctx context.Context, clientID string) error
	StoreClientData(ctx context.Context, clientID string, cd *ClientData) error
}

// MysqlDB ...
type MysqlDB struct {
	Dsn string
	Cxn *sql.DB
	Log *logrus.Logger
}

// NewMysqlDB creates a new DatabaseIF for mysql db
func NewMysqlDB(dataSourseName string, log *logrus.Logger) (DatabaseIF, error) {

	mysql := &MysqlDB{
		Dsn: dataSourseName,
		Log: log,
	}
	if err := mysql.Connect(); err != nil {
		return nil, err
	}

	return mysql, nil
}

// Connect connects to the DB instance
func (mysql *MysqlDB) Connect() error {

	cxn, err := sql.Open(mysqlDriver, mysql.Dsn)
	if err != nil {
		return errors.Wrap(err, "could not create connection to MYSQL db")
	}
	mysql.Cxn = cxn
	if err := mysql.Test(); err != nil {
		return err
	}

	return nil
}

// Test test the DB connection to see if connection if necessary
func (mysql *MysqlDB) Test() error {

	if err := mysql.Cxn.Ping(); err != nil {
		return errors.Wrap(err, "could not conect to MYSQL db")
	}
	return nil
}

// PrepareTableForNewClient creates a new table for the new client
func (mysql *MysqlDB) PrepareTableForNewClient(ctx context.Context, clientID string) error {

	stmt := fmt.Sprintf("CREATE TABLE `%v` (`id` int(11) NOT NULL AUTO_INCREMENT,`createdAt` timestamp NULL DEFAULT NULL,`expiresAt` timestamp NULL DEFAULT NULL,`data` varchar(500) DEFAULT NULL,PRIMARY KEY (`id`))", clientID)

	_, err := mysql.Cxn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

// GenerateClientID gets clientId from clients table
func (mysql *MysqlDB) GenerateClientID(ctx context.Context, name string) (string, error) {

	clientID := uuid.New().String()

	stmt := `INSERT INTO Clients (clientName,clientId) VALUES(?,?)`

	_, err := mysql.Cxn.ExecContext(ctx, stmt, name, clientID)
	if err != nil {
		return "", err
	}

	return clientID, nil
}

// GetClientID gets clientName based on the given clientId
func (mysql *MysqlDB) GetClientID(ctx context.Context, name string) (string, error) {

	var clientID string

	stmt := `SELECT clientId FROM Clients where clientName = ?`
	err := mysql.Cxn.QueryRowContext(ctx, stmt, name).Scan(&clientID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return clientID, nil
}

// StoreClientData persists client data based on the given clientID
func (mysql *MysqlDB) StoreClientData(ctx context.Context, clientID string, cd *ClientData) error {

	stmt := fmt.Sprintf("INSERT INTO `%v` (createdAt, expiresAt, data) VALUES(%v,%v,%v)", clientID, cd.CreatedAt, cd.ExpiresAt, cd.Data)

	_, err := mysql.Cxn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}
