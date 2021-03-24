package syncman

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

	db, err := sql.Open("postgres", connectionstring)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Error connecting to database"), err, nil)

	}
	defer func() { _ = db.Close() }()

	if err = db.Ping(); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Ping to the database was not successful"), err, nil)
	}

	sqlstatement := `CREATE SCHEMA IF NOT EXISTS` + dbschemaname
	// Exec executes a query without returning any rows.
	if _, err = db.Exec(sqlstatement); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("schema creation query failed (%s)", sqlstatement), err, nil)
	}

	configResource := `CREATE TABLE IF NOT EXISTS` + dbschemaname + `.` + scConfig + `( resource_id text PRIMARY KEY, resource_type text, resource text, project text, cluster_id text)`

	// Exec executes a query without returning any rows.
	if _, err = db.Exec(configResource); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Table creation query failed (%s)", configResource), err, nil)
	}
	return &PostgresStore{db: db, globalConfig: config.GenerateEmptyConfig(), clusterID: clusterID, dbschemaname: dbschemaname}, nil
}

// Register registers space cloud to the postgres store
func (s *PostgresStore) Register() {}

// WatchResources maintains consistency over all projects
func (s *PostgresStore) WatchResources(cb func(eventType, resourceId string, resourceType config.Resource, resource interface{})) error {
	type configResource struct {
		resourceID   string      `db:"resource_id"`
		resourceType interface{} `db:"resource_type"`
		resource     string      `db:"resource"`
		project      string      `db:"project"`
		clusterID    string      `db:"cluster_id"`
	}
	var r configResource
	configResourceMap := make(map[string]postgres)
	var res interface{}

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	go func() {
		for range t.C {

			rows, err := s.db.Query("SELECT resource_id, resource_type, resource, project FROM $1 WHERE cluster_id = $2", fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), s.clusterID)
			if err != nil {
				continue
			}

			for rows.Next() {
				if err := rows.Scan(&r); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to get rows from database", err, nil)
					break
				}
				if err := json.Unmarshal([]byte(r.resource), &res); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to unmarshal resource from database", err, nil)
					break
				}
				configResourceMap[r.resourceID] = postgres{
					resourceID:   r.resourceID,
					resourceType: r.resourceType,
					resource:     res,
					project:      r.project,
					clusterID:    r.clusterID,
				}
			}
			if err = rows.Err(); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error was encountered during iteration of rows in database", err, nil)
				break
			}

			for _, i := range configResourceMap {
				evenType, resourceID, resourceType, resource := s.helperAddOrUpdate(i.resourceID, i.resourceType.(config.Resource), i.resource, i.project)
				if resource == nil || resourceID == "" {
					break
				}
				cb(evenType, resourceID, resourceType, resource)
			}

			// Delete
			var obj interface{}
			for project, configs := range s.globalConfig.Projects {

				// project
				projectconfig := configs.ProjectConfig
				resID := config.GenerateResourceID(s.clusterID, projectconfig.ID, config.ResourceProject, projectconfig.ID)
				if _, ok := configResourceMap[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Database Config
				databaseConfigs := configs.DatabaseConfigs
				for resourceID := range databaseConfigs {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Database rule
				databaseRule := configs.DatabaseRules
				for resourceID := range databaseRule {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Database schema
				databaseSchema := configs.DatabaseSchemas
				for resourceID := range databaseSchema {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				//Database Prepared Queries
				databasePrepQueries := configs.DatabasePreparedQueries
				for resourceID := range databasePrepQueries {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing config
				// eventingConfig := configs.EventingConfig
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceEventingConfig, "eventing")
				if _, ok := configResourceMap[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Eventing Schema
				eventingSchema := configs.EventingSchemas
				for resourceID := range eventingSchema {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing Rule
				eventingRule := configs.EventingRules
				for resourceID := range eventingRule {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Eventing trigger
				eventingTrigger := configs.EventingTriggers
				for resourceID := range eventingTrigger {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// FileStoreConfig
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceFileStoreConfig, "filestore")
				if _, ok := configResourceMap[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// FileStoreRule
				fileStoreRule := configs.FileStoreRules
				for resourceID := range fileStoreRule {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Auths
				auths := configs.Auths
				for resourceID := range auths {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// LetsEncrypt
				LetsEncrypt := configs.LetsEncrypt
				if _, ok := configResourceMap[LetsEncrypt.ID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Ingress Routes
				ingresRoutes := configs.IngressRoutes
				for resourceID := range ingresRoutes {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}

				// Ingress Global
				// ingressGlobal := configs.IngressGlobal
				resID = config.GenerateResourceID(s.clusterID, project, config.ResourceIngressGlobal, "global")
				if _, ok := configResourceMap[resID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}

				// Service
				services := configs.RemoteService
				for resourceID := range services {
					if _, ok := configResourceMap[resourceID]; !ok {
						evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
						if resource == nil || resourceID == "" {
							break
						}
						cb(evenType, resourceID, resourceType, resource)
					}
				}
			}

			// Cluster
			// Cluster :=s.globalConfig.ClusterConfig
			resID := config.GenerateResourceID(s.clusterID, "", config.ResourceCluster, "cluster")
			if _, ok := configResourceMap[resID]; !ok {
				evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
				if resource == nil || resourceID == "" {
					break
				}
				cb(evenType, resourceID, resourceType, resource)
			}

			// Integration
			for resourceID := range s.globalConfig.Integrations {
				if _, ok := configResourceMap[resourceID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}
			}

			// IntegrationHook
			for resourceID := range s.globalConfig.IntegrationHooks {
				if _, ok := configResourceMap[resourceID]; !ok {
					evenType, resourceID, resourceType, resource := onAddOrUpdateResource(config.ResourceDeleteEvent, obj)
					if resource == nil || resourceID == "" {
						break
					}
					cb(evenType, resourceID, resourceType, resource)
				}
			}
			_ = rows.Close()
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

	query := "SELECT * FROM $1 WHERE resource_id=$2 AND project=$3 AND cluster_id=$4"
	row := s.db.QueryRow(query, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID, projectID, clusterID)
	var count int
	err = row.Scan(&count)

	switch err {
	case sql.ErrNoRows:
		sqlStatement := `INSERT INTO $1 (resource_id, resource_type, resource, project, cluster_id) VALUES ($2, $3, $4, $5, $6)`
		_, err = s.db.Exec(sqlStatement, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID, resourceType, string(res), projectID, clusterID)
		if err != nil {
			return err
		}
		return nil
	case nil:
		sqlStatement := `UPDATE $1 SET resource_type = $3, resource = $4 WHERE resource_id = $2 AND project=$5 AND cluster_id=$6;`
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

	clusterID, projectID, _, err := splitResourceID(ctx, resourceID)
	if err != nil {
		return err
	}

	sqlStatement := `DELETE FROM $1 WHERE resource_id = $2 AND project=$3 AND cluster_id=$4;`
	_, err = s.db.Exec(sqlStatement, fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceID, projectID, clusterID)
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

		rows, err := s.db.Query("SELECT resource_id,resource FROM $1 WHERE resource_type = $2 AND cluster_id = $3", fmt.Sprintf("%s.%s", s.dbschemaname, scConfig), resourceType, s.clusterID)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			if err := rows.Scan(&resourceID, &res); err != nil {
				return nil, err
			}
			if err := json.Unmarshal([]byte(res), &resource); err != nil {
				return nil, err
			}
			if err := updateResource(context.TODO(), config.ResourceAddEvent, globalConfig, resourceID, resourceType, resource); err != nil {
				return nil, err
			}

			_ = rows.Close()
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		_ = rows.Close()
	}
	return globalConfig, nil
}
