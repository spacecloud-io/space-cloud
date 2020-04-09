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

// RawExec performs an operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawExec(ctx context.Context, query string) error {
	_, err := s.client.ExecContext(ctx, query, []interface{}{}...)
	return err
}

// RawQuery query document(s) from the database
func (s *SQL) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error) {
	return s.readexec(ctx, query, args, utils.All, s.client)
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
func (s *SQL) CreateDatabaseIfNotExist(ctx context.Context, project string) error {
	var sql string
	switch utils.DBType(s.dbType) {
	case utils.MySQL:
		sql = "create database if not exists " + project
	case utils.Postgres:
		sql = "create schema if not exists " + project
	case utils.SQLServer:
		sql = `IF (NOT EXISTS (SELECT * FROM sys.schemas WHERE name = '` + project + `')) 
					BEGIN
    					EXEC ('CREATE SCHEMA [` + project + `] ')
					END`
	default:
		return fmt.Errorf("invalid db type (%s) provided", s.dbType)
	}
	return s.RawExec(ctx, sql)
}
