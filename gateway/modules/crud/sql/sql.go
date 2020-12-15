package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spaceuptech/helpers"

	_ "github.com/denisenkom/go-mssqldb" // Import for MsSQL
	_ "github.com/go-sql-driver/mysql"   // Import for MySQL
	_ "github.com/lib/pq"                // Import for postgres

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SQL holds the sql db object
type SQL struct {
	lock                sync.RWMutex
	enabled             bool
	queryFetchLimit     *int64
	connection          string
	client              *sqlx.DB
	dbType              string
	name                string // logical db name or schema name according to the database type
	auth                model.AuthCrudInterface
	driverConf          config.DriverConfig
	connRetryCloserChan chan struct{}
}

// Init initialises a new sql instance
func Init(dbType model.DBType, enabled bool, connection string, dbName string, auth model.AuthCrudInterface, driverConf config.DriverConfig) (s *SQL, err error) {
	s = &SQL{enabled: enabled, connection: connection, name: dbName, client: nil, auth: auth, driverConf: driverConf}

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

	closer := make(chan struct{}, 1)
	s.connRetryCloserChan = closer
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				if !s.GetConnectionState(ctx) {
					if err := s.connect(); err != nil {
						_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Automatic connection retry failed for (%s) db with logical db name (%s)", dbType, dbName), err, nil)
					}
				}
				cancel()
			case <-closer:
				close(closer)
				ticker.Stop()
				return
			}
		}
	}()
	return
}

// IsSame checks if we've got the same connection string
func (s *SQL) IsSame(conn, dbName string, driverConf config.DriverConfig) bool {
	return strings.HasPrefix(s.connection, conn) && dbName == s.name && driverConf.MaxConn == s.driverConf.MaxConn && driverConf.MaxIdleTimeout == s.driverConf.MaxIdleTimeout && driverConf.MaxIdleConn == s.driverConf.MaxIdleConn
}

// Close gracefully the SQL client
func (s *SQL) Close() error {
	if s.getClient() != nil {
		s.connRetryCloserChan <- struct{}{}
		if err := s.getClient().Close(); err != nil {
			_ = helpers.Logger.LogError("close", fmt.Sprintf("Unable to close (%s) db (%s) connection", s.dbType, s.name), err, nil)
		}
		s.setClient(nil)
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
	if s.getClient() == nil {
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

	s.setClient(sql)

	maxConn := s.driverConf.MaxConn
	if maxConn == 0 {
		maxConn = 100
	}

	maxIdleConn := s.driverConf.MaxIdleConn
	if maxIdleConn == 0 {
		maxIdleConn = 50
	}

	maxIdleTimeout := s.driverConf.MaxIdleTimeout
	if maxIdleTimeout == 0 {
		maxIdleTimeout = 60 * 5 * 1000
	}

	s.getClient().SetMaxOpenConns(maxConn)
	s.getClient().SetMaxIdleConns(maxIdleConn)
	duration := time.Duration(maxIdleTimeout) * time.Millisecond
	s.getClient().SetConnMaxIdleTime(duration)
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

// SetQueryFetchLimit sets data fetch limit
func (s *SQL) SetQueryFetchLimit(limit int64) {
	s.queryFetchLimit = &limit
}

func (s *SQL) setClient(c *sqlx.DB) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.client = c
}

func (s *SQL) getClient() *sqlx.DB {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.client
}
