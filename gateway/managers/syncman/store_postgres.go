package syncman

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

// NewPostgresStore creates a new Postgres Store
func NewPostgresStore(connectionstring, clusterID string) (*PostgresStore, error) {
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
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	const qry = `
	CREATE TABLE IF NOT EXISTS SC_CONFIG (
		resourceId text PRIMARY KEY,
		resourceType text,
		resource text,
		project text,
		clusterId text
	)`

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
	var obj interface{}
	go func() {
		for range time.Tick(10 * time.Second) {
			rows, err := s.db.Query("select resourceId,resourceType,resource,project from users where clusterID = $2", s.clusterID)
			if err != nil {
				return
			}
			defer rows.Close()
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
				switch resourceType.(config.Resource) {
				case config.ResourceProject:
					projectConfig := s.globalConfig.Projects[project].ProjectConfig
					if !reflect.DeepEqual(projectConfig, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceDatabaseConfig:
					databaseConfig, ok := s.globalConfig.Projects[project].DatabaseConfigs[resourceID]
					if ok {
						if !reflect.DeepEqual(databaseConfig, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceDatabaseRule:
					databaseRule, ok := s.globalConfig.Projects[project].DatabaseRules[resourceID]
					if ok {
						if !reflect.DeepEqual(databaseRule, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceDatabaseSchema:
					databaseSchema, ok := s.globalConfig.Projects[project].DatabaseSchemas[resourceID]
					if ok {
						if !reflect.DeepEqual(databaseSchema, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceDatabasePreparedQuery:
					databasePreparedQuery, ok := s.globalConfig.Projects[project].DatabasePreparedQueries[resourceID]
					if ok {
						if !reflect.DeepEqual(databasePreparedQuery, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceFileStoreConfig:
					fileStoreConfig := s.globalConfig.Projects[project].FileStoreConfig
					if !reflect.DeepEqual(fileStoreConfig, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceFileStoreRule:
					fileStoreRule, ok := s.globalConfig.Projects[project].FileStoreRules[resourceID]
					if ok {
						if !reflect.DeepEqual(fileStoreRule, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceEventingConfig:
					eventingConfig := s.globalConfig.Projects[project].EventingConfig
					if !reflect.DeepEqual(eventingConfig, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceEventingTrigger:
					eventingTrigger, ok := s.globalConfig.Projects[project].EventingTriggers[resourceID]
					if ok {
						if !reflect.DeepEqual(eventingTrigger, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceEventingRule:
					eventingRule, ok := s.globalConfig.Projects[project].EventingRules[resourceID]
					if ok {
						if !reflect.DeepEqual(eventingRule, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceEventingSchema:
					eventingSchemas, ok := s.globalConfig.Projects[project].EventingSchemas[resourceID]
					if ok {
						if !reflect.DeepEqual(eventingSchemas, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceRemoteService:
					remoteService, ok := s.globalConfig.Projects[project].RemoteService[resourceID]
					if ok {
						if !reflect.DeepEqual(remoteService, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceIngressGlobal:
					global := s.globalConfig.Projects[project].IngressGlobal
					if !reflect.DeepEqual(global, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceIngressRoute:
					routes, ok := s.globalConfig.Projects[project].IngressRoutes[resourceID]
					if ok {
						if !reflect.DeepEqual(routes, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceAuthProvider:
					auth, ok := s.globalConfig.Projects[project].Auths[resourceID]
					if ok {
						if !reflect.DeepEqual(auth, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceProjectLetsEncrypt:
					letsEncrypt := s.globalConfig.Projects[project].LetsEncrypt
					if !reflect.DeepEqual(letsEncrypt, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceCluster:
					cluster := s.globalConfig.ClusterConfig
					if !reflect.DeepEqual(cluster, resource) {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceIntegration:
					integrationconfig, ok := s.globalConfig.Integrations[resourceID]
					if ok {
						if !reflect.DeepEqual(integrationconfig, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				case config.ResourceIntegrationHook:
					integrationHook, ok := s.globalConfig.IntegrationHooks[resourceID]
					if ok {
						if !reflect.DeepEqual(integrationHook, resource) {
							evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceUpdateEvent, obj)
							if resource == nil || resourceID == "" {
								return
							}
							cb(evenType, resourceID, resourceType, resource)
						}
					} else {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceAddEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}
			}
			err = rows.Err()
			if err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error was encountered during iteration of rows in database", err, nil)
				return
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
	res, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	query := "SELECT COUNT(*) as count FROM SC_CONFIG where resourceid=$1"
	row := s.db.QueryRow(query, resourceID)
	var count int
	err = row.Scan(&count)

	switch err {
	case sql.ErrNoRows:
		sqlStatement := `
INSERT INTO users (resourceId, resourceType, resource, project, clusterId)
VALUES ($1, $2, $3, $4, $5)`
		_, err = s.db.Exec(sqlStatement, resourceID, resourceType, string(res), projectID, clusterID)
		if err != nil {
			return err
		}
		return nil
	case nil:
		sqlStatement := `
UPDATE users
SET resourceType = $2, resource = $3, project = $4, clusterId = $5
WHERE resourceId = $1;`
		res, err := s.db.Exec(sqlStatement, resourceID, resourceType, string(res), projectID, clusterID)
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
	sqlStatement := `
DELETE FROM users
WHERE resourceId = $1;`
	_, err := s.db.Exec(sqlStatement, resourceID)
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
		rows, err := s.db.Query("select resourceId,resource from users where resourceType = $1 and clusterID = $2", resourceType, s.clusterID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&resourceID, &res)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal([]byte(res), &resource); err != nil {
				return nil, err
			}
			if err := updateResource(context.TODO(), "add", globalConfig, resourceID, resourceType, resource); err != nil {
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
