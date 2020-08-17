package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawBatch(ctx context.Context, queries []string) error {
	// Skip if length of queries == 0
	if len(queries) == 0 {
		return nil
	}

	logrus.Debugf("Executing sql raw query - %v", queries)

	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// RawQuery query document(s) from the database
func (s *SQL) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error) {
	return s.readexec(ctx, query, args, utils.All, s.client, false)
}

// GetConnectionState : Function to get connection state
func (s *SQL) GetConnectionState(ctx context.Context) bool {
	if !s.enabled || s.client == nil {
		return false
	}

	// Ping to check if connection is established
	err := s.client.PingContext(ctx)
	return err == nil
}

// CreateDatabaseIfNotExist creates a schema / database
func (s *SQL) CreateDatabaseIfNotExist(ctx context.Context, name string) error {
	var sql string
	switch utils.DBType(s.dbType) {
	case utils.MySQL:
		sql = "create database if not exists " + name
	case utils.Postgres:
		sql = "create schema if not exists " + name
	case utils.SQLServer:
		sql = `IF (NOT EXISTS (SELECT * FROM sys.schemas WHERE name = '` + name + `')) 
					BEGIN
    					EXEC ('CREATE SCHEMA [` + name + `]')
					END`
	default:
		return fmt.Errorf("invalid db type (%s) provided", s.dbType)
	}
	return s.RawBatch(ctx, []string{sql})
}
