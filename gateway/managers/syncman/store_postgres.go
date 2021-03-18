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
	dbschemaname string
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
	return &PostgresStore{db: db, globalConfig: conf, clusterID: clusterID, dbschemaname: dbschemaname}, nil
}

// Register registers space cloud to the postgres store
func (s *PostgresStore) Register() {}

// WatchResources maintains consistency over all projects
func (s *PostgresStore) WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error {
	type qry struct {
		resourceID   string      `db:"resourceId"`
		resourceType interface{} `db:"resourceType"`
		resource     string      `db:"resource"`
		project      string      `db:"project"`
		clusterID    string      `db:"clusterId"`
	}
	var r qry
	resources := make(map[string]postgres)
	var res interface{}

	var rows *sql.Rows
	go func() {
		for range time.NewTicker(10 * time.Second).C {

			rows, err := s.db.Query("SELECT resourceId,resourceType,resource,project FROM $1 WHERE clusterID = $2", fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), s.clusterID)
			if err != nil {
				return
			}

			for rows.Next() {
				err := rows.Scan(&r)
				if err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to get rows from database", err, nil)
					return
				}
				if err := json.Unmarshal([]byte(r.resource), &res); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource from database", err, nil)
					return
				}
				resources[r.resourceID] = postgres{
					resourceID:   r.resourceID,
					resourceType: r.resourceType,
					resource:     res,
					project:      r.project,
					clusterID:    r.clusterID,
				}
			}
			err = rows.Err()
			if err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error was encountered during iteration of rows in database", err, nil)
				return
			}

			for _, i := range resources {
				evenType, resourceID, resourceType, resource := s.helperAddOrUpdate(i.resourceID, i.resourceType.(config.Resource), i.resource, i.project)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			}

			// Delete
			var obj interface{}
			for project, configs := range s.globalConfig.Projects {

				// project
				projectconfig := configs.ProjectConfig
				resID := config.GenerateResourceID(s.clusterID, projectconfig.ID, config.ResourceProject, projectconfig.ID)
				if _, ok := resources[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Database Config
				databaseConfigs := configs.DatabaseConfigs
				for resourceID := range databaseConfigs {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Database rule
				databaseRule := configs.DatabaseRules
				for resourceID := range databaseRule {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Database schema
				databaseSchema := configs.DatabaseSchemas
				for resourceID := range databaseSchema {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				//Database Prepared Queries
				databasePrepQueries := configs.DatabasePreparedQueries
				for resourceID := range databasePrepQueries {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing config
				// eventingConfig := configs.EventingConfig
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceEventingConfig, "eventing")
				if _, ok := resources[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Eventing Schema
				eventingSchema := configs.EventingSchemas
				for resourceID := range eventingSchema {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing Rule
				eventingRule := configs.EventingRules
				for resourceID := range eventingRule {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing trigger
				eventingTrigger := configs.EventingTriggers
				for resourceID := range eventingTrigger {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// FileStoreConfig
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreConfig, "filestore")
				if _, ok := resources[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// FileStoreRule
				fileStoreRule := configs.FileStoreRules
				for resourceID := range fileStoreRule {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Auths
				auths := configs.Auths
				for resourceID := range auths {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// LetsEncrypt
				LetsEncrypt := configs.LetsEncrypt
				if _, ok := resources[LetsEncrypt.ID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Ingress Routes
				ingresRoutes := configs.IngressRoutes
				for resourceID := range ingresRoutes {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Ingress Global
				// ingressGlobal := configs.IngressGlobal
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceIngressGlobal, "global")
				if _, ok := resources[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Service
				services := configs.RemoteService
				for resourceID := range services {
					if _, ok := resources[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							return
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}
			}

			// Cluster
			// Cluster :=s.globalConfig.ClusterConfig
			resID := config.GenerateResourceID(s.clusterID, "", config.ResourceCluster, "cluster")
			if _, ok := resources[resID]; !ok {
				evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
				if resource == nil || resourceID == "" {
					return
				}
				cb(evenType, resourceID, resourceType, resource)
			}

			// Integration
			for resourceID := range s.globalConfig.Integrations {
				if _, ok := resources[resourceID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}
			}

			// IntegrationHook
			for resourceID := range s.globalConfig.IntegrationHooks {
				if _, ok := resources[resourceID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						return
					}
					cb(evenType, resourceID, resourceType, resource)
				}
			}

		}
	}()
	defer func() { _ = rows.Close() }()
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
	row := s.db.QueryRow(query, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID)
	var count int
	err = row.Scan(&count)

	switch err {
	case sql.ErrNoRows:
		sqlStatement := `INSERT INTO $1 (resourceId, resourceType, resource, project, clusterId) VALUES ($2, $3, $4, $5, $6)`
		_, err = s.db.Exec(sqlStatement, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID, resourceType, string(res), projectID, clusterID)
		if err != nil {
			return err
		}
		return nil
	case nil:
		sqlStatement := `UPDATE $1 SET resourceType = $3, resource = $4, project = $5, clusterId = $6 WHERE resourceId = $2;`
		res, err := s.db.Exec(sqlStatement, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID, resourceType, string(res), projectID, clusterID)
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
	_, err := s.db.Exec(sqlStatement, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID)
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
		rows, err := s.db.Query("SELECT resourceId,resource FROM $1 WHERE resourceType = $2 AND clusterID = $3", fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceType, s.clusterID)
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
