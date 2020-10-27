package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spaceuptech/helpers"

	_ "github.com/denisenkom/go-mssqldb" // Import for MsSQL
	_ "github.com/go-sql-driver/mysql"   // Import for MySQL
	_ "github.com/lib/pq"                // Import for postgres

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SQL holds the sql db object
type SQL struct {
	enabled    bool
	connection string
	client     *sqlx.DB
	dbType     string
	name       string // logical db name or schema name according to the database type
	auth       model.AuthCrudInterface
}

// Init initialises a new sql instance
func Init(dbType model.DBType, enabled bool, connection string, dbName string, auth model.AuthCrudInterface) (s *SQL, err error) {
	s = &SQL{enabled: enabled, connection: connection, name: dbName, client: nil, auth: auth}

	switch dbType {
	case model.Postgres:
		s.dbType = "postgres"

	case model.MySQL:
		s.dbType = "mysql"

	case model.SQLServer:
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

// IsSame checks if we've got the same connection string
func (s *SQL) IsSame(conn, dbName string) bool {
	return strings.HasPrefix(s.connection, conn) && dbName == s.name
}

// Close gracefully the SQL client
func (s *SQL) Close() error {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			return err
		}

		s.client = nil
	}

	return nil
}

// GetDBType returns the dbType of the crud block
func (s *SQL) GetDBType() model.DBType {
	switch s.dbType {
	case "postgres":
		return model.Postgres
	case "mysql":
		return model.MySQL
	case "sqlserver":
		return model.SQLServer
	}

	return model.MySQL
}

// IsClientSafe checks whether database is enabled and connected
func (s *SQL) IsClientSafe(ctx context.Context) error {
	if !s.enabled {
		return utils.ErrDatabaseDisabled
	}

	if s.client == nil {
		if err := s.connect(); err != nil {
			helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Error connecting to "+s.dbType+" : "+err.Error()), nil)
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

	s.client.SetMaxOpenConns(10)
	s.client.SetMaxIdleConns(0)
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
	defer func() { _ = stmt.Close() }()

	return stmt.ExecContext(ctx, args...)
}
