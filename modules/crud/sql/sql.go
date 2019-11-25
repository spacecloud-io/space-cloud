package sql

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/denisenkom/go-mssqldb" //Import for MsSQL
	_ "github.com/go-sql-driver/mysql"   // Import for MySQL
	_ "github.com/lib/pq"                // Import for postgres

	"github.com/spaceuptech/space-cloud/utils"
)

// SQL holds the sql db object
type SQL struct {
	enabled            bool
	connection         string
	client             *sqlx.DB
	dbType             string
	removeProjectScope bool
}

// Init initialises a new sql instance
func Init(dbType utils.DBType, enabled, removeProjectScope bool, connection string) (s *SQL, err error) {
	s = &SQL{enabled: enabled, removeProjectScope: removeProjectScope, connection: connection, client: nil}

	switch dbType {
	case utils.Postgres:
		s.dbType = "postgres"

	case utils.MySQL:
		s.dbType = "mysql"

	case utils.SqlServer:
		s.dbType = "sqlserver"

	default:
		err = utils.ErrUnsupportedDatabase
		return
	}

	if s.enabled {
		err = s.connect()
	}

	return
}

// Close gracefully the SQL client
func (s *SQL) Close() error {
	if s.client != nil {
		return s.client.Close()
	}

	return nil
}

// GetDBType returns the dbType of the crud block
func (s *SQL) GetDBType() utils.DBType {
	switch s.dbType {
	case "postgres":
		return utils.Postgres
	case "mysql":
		return utils.MySQL
	case "sqlserver":
		return utils.SqlServer
	}

	return utils.MySQL
}

// IsClientSafe checks whether database is enabled and connected
func (s *SQL) IsClientSafe() error {
	if !s.enabled {
		return utils.ErrDatabaseDisabled
	}

	if s.client == nil {
		if err := s.connect(); err != nil {
			log.Println("Error connecting to " + s.dbType + " : " + err.Error())
			return utils.ErrDatabaseConnection
		}
	}

	return nil
}

func (s *SQL) connect() error {
	timeOut := 3 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	sql, err := sqlx.Open(s.dbType, s.connection)
	if err != nil {
		return err
	}

	s.client = sql

	return sql.PingContext(ctx)
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
