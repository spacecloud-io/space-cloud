package sql

import (
	"context"
	"errors"
	"time"
	"database/sql"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" // Import for MySQL
	_ "github.com/lib/pq"              // Import for postgres

	"github.com/spaceuptech/space-cloud/utils"
)

// SQL holds the sql db object
type SQL struct {
	client  *sqlx.DB
	dbType  string
	timeOut time.Duration
}

// Init initialises a new sql instance
func Init(dbType utils.DBType, connection string) (*SQL, error) {
	var sql *sqlx.DB
	var err error

	s := &SQL{}

	timeOut := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	switch dbType {
	case utils.Postgres:
		sql, err = sqlx.Open("postgres", connection)
		s.dbType = "postgres"

	case utils.MySQL:
		sql, err = sqlx.Open("mysql", connection)
		s.dbType = "mysql"

	default:
		return nil, errors.New("SQL: Invalid driver provided")
	}

	if err != nil {
		return nil, err
	}

	err = sql.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	s.client = sql
	s.timeOut = timeOut
	return s, nil
}

// Close gracefully the SQL client
func (s *SQL) Close() error {
	return s.client.Close()
}

// GetDBType returns the dbType of the crud block
func (s *SQL) GetDBType() utils.DBType {
	switch s.dbType {
	case "postgres":
		return utils.Postgres
	case "mysql":
		return utils.MySQL
	}

	return utils.MySQL
}

type executor interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

func doExecContext(ctx context.Context, query string, args []interface{}, executor executor) (sql.Result, error) {
	stmt, err := executor.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}
