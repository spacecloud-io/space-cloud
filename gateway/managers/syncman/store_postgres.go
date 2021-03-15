package syncman

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// PostgresStore is an object for storing PostgresStore information
type PostgresStore struct {
	db           *sql.DB
	globalConfig *config.Config
	clusterID    string
}

type postgres struct {
	resourceID   string
	resourceType interface{}
	resource     interface{}
	project      string
	clusterID    string
}

const scConfig string = "SC_CONFIG"

// NewPostgresStore creates a new Postgres Store
func NewPostgresStore(connectionstring, clusterID, dbschemaname string) (*PostgresStore, error) {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}
	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}

	// For compatibility with v18
	if conf.ClusterConfig == nil {
		conf.ClusterConfig = &config.ClusterConfig{EnableTelemetry: true}
	}

	db, err := sql.Open("postgres", connectionstring)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Error connecting to database"), err, nil)

	}
	defer func() { _ = db.Close() }()

	err = db.Ping()
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Ping to the database was not successful"), err, nil)
	}

	sqlstatement := `CREATE SCHEMA IF NOT EXISTS` + dbschemaname
	// Exec executes a query without returning any rows.
	if _, err = db.Exec(sqlstatement); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("schema creation query failed (%s)", sqlstatement), err, nil)
	}

	qry := `CREATE TABLE IF NOT EXISTS` + dbschemaname + `.` + scConfig + `( resourceId text PRIMARY KEY, resourceType text, resource text, project text, clusterId text)`

	// Exec executes a query without returning any rows.
	if _, err = db.Exec(qry); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Table creation query failed (%s)", qry), err, nil)
	}
	return &PostgresStore{db: db, globalConfig: conf, clusterID: clusterID}, nil
}

// Register registers space cloud to the postgres store
func (s *PostgresStore) Register() {}

// WatchResources maintains consistency over all projects
func (s *PostgresStore) WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error {
	var (
		resourceID   string
		resource     interface{}
		res          string
		resourceType interface{}
		project      string
	)

	resources := make([]postgres, 0)
	go func() {
		for range time.NewTicker(10 * time.Second).C {

			rows, err := s.db.Query("SELECT resourceId,resourceType,resource,project FROM $1 WHERE clusterID = $2", scConfig, s.clusterID)
			if err != nil {
				return
			}
			defer func() { _ = rows.Close() }()

			for rows.Next() {
				err := rows.Scan(&resourceID, &resourceType, &res, &project)
				if err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to get rows from database", err, nil)
					return
				}
				if err := json.Unmarshal([]byte(res), &resource); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource from database", err, nil)
					return
				}
				resources = append(resources, postgres{resourceID: resourceID, resourceType: resourceType, resource: resource, project: project, clusterID: s.clusterID})
			}
			err = rows.Err()
			if err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error was encountered during iteration of rows in database", err, nil)
				return
			}

			for _, i := range resources {
				evenType, resourceID, resourceType, resource := s.helper(resources, i.resourceID, i.resourceType.(config.Resource), i.resource, i.project)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			}
		}
	}()
	return nil
}

// SetResource sets the project of the postgres globalConfig
func (s *PostgresStore) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	if err := updateResource(ctx, config.ResourceAddEvent, s.globalConfig, resourceID, "", resource); err != nil {
		return err
	}
	clusterID, projectID, resourceType, err := splitResourceID(ctx, resourceID)
	if err != nil {
		return err
	}

	res, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	query := "SELECT * FROM $1 WHERE resourceid=$2"
	row := s.db.QueryRow(query, scConfig, resourceID)
	var count int
	err = row.Scan(&count)

	switch err {
	case sql.ErrNoRows:
		sqlStatement := `INSERT INTO $1 (resourceId, resourceType, resource, project, clusterId) VALUES ($2, $3, $4, $5, $6)`
		_, err = s.db.Exec(sqlStatement, scConfig, resourceID, resourceType, string(res), projectID, clusterID)
		if err != nil {
			return err
		}
		return nil
	case nil:
		sqlStatement := `UPDATE $1 SET resourceType = $3, resource = $4, project = $5, clusterId = $6 WHERE resourceId = $2;`
		res, err := s.db.Exec(sqlStatement, scConfig, resourceID, resourceType, string(res), projectID, clusterID)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// DeleteResource deletes the project from the postgres gloablConfig
func (s *PostgresStore) DeleteResource(ctx context.Context, resourceID string) error {
	if err := updateResource(ctx, config.ResourceDeleteEvent, s.globalConfig, resourceID, "", nil); err != nil {
		return err
	}
	sqlStatement := `DELETE FROM $1 WHERE resourceId = $2;`
	_, err := s.db.Exec(sqlStatement, scConfig, resourceID)
	return err
}

// GetGlobalConfig gets config all projects
func (s *PostgresStore) GetGlobalConfig() (*config.Config, error) {
	globalConfig := config.GenerateEmptyConfig()
	var (
		resourceID string
		resource   interface{}
		res        string
	)
	for _, resourceType := range config.ResourceFetchingOrder {
		rows, err := s.db.Query("SELECT resourceId,resource FROM $1 WHERE resourceType = $2 AND clusterID = $3", scConfig, resourceType, s.clusterID)
		if err != nil {
			return nil, err
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			err := rows.Scan(&resourceID, &res)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal([]byte(res), &resource); err != nil {
				return nil, err
			}
			if err := updateResource(context.TODO(), config.ResourceAddEvent, globalConfig, resourceID, resourceType, resource); err != nil {
				return nil, err
			}
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}
	return globalConfig, nil
}
