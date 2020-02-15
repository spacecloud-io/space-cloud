package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleGetCollections is an endpoint handler which return all the collection(table) names for specified data base
func HandleGetCollections(adminMan *admin.Manager, crud *crud.Module, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		collections, err := crud.GetCollections(ctx, project, dbType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		cols := make([]string, len(collections))
		for i, value := range collections {
			cols[i] = value.TableName
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"collections": cols})
	}
}

// HandleGetConnectionState gives the status of connection state of client
func HandleGetConnectionState(adminMan *admin.Manager, crud *crud.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbType := vars["dbType"]

		connState := crud.GetConnectionState(ctx, dbType)

		w.WriteHeader(http.StatusOK) // http status code
		_ = json.NewEncoder(w).Encode(map[string]bool{"status": connState})
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database & removes it from config
func HandleDeleteCollection(adminMan *admin.Manager, crud *crud.Module, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		if err := crud.DeleteTable(ctx, project, dbType, col); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetDeleteCollection(ctx, project, dbType, col); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status code
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleDatabaseConnection is an endpoint handler which updates database config & connects to database
func HandleDatabaseConnection(adminMan *admin.Manager, crud *crud.Module, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.CrudStub{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]

		if err := syncman.SetDatabaseConnection(ctx, project, dbType, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

//HandleGetDatabaseConnection returns handler to get Database Collection
func HandleGetDatabaseConnection(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbType, exists := r.URL.Query()["dbType"]
		if !exists {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("dbType not provided in url")})
			return
		}

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get collection
		coll, ok := project.Modules.Crud[dbType[0]]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("collection not found", r.URL)})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"conn": coll.Conn, "enabled": coll.Enabled, "collections": coll.Collections})
	}
}

// HandleRemoveDatabaseConfig is an endpoint handler which removes database config
func HandleRemoveDatabaseConfig(adminMan *admin.Manager, crud *crud.Module, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbType"]
		project := vars["project"]

		if err := syncman.RemoveDatabaseConfig(ctx, project, dbAlias); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}

// HandleModifySchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifySchema(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if err := schemaArg.SchemaModifyAll(ctx, dbType, project, map[string]*config.TableRule{col: &v}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetModifySchema(ctx, project, dbType, col, v.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}

//HandleGetSchema returns handler to get schema
func HandleGetSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		dbType, exists := r.URL.Query()["dbType"]
		if !exists {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("dbType not provided in url")})
		}

		//gel col from url
		col, exists := r.URL.Query()["col"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get collection
		coll, ok := project.Modules.Crud[dbType[0]]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("collection not found")})
			return
		}

		//check if col exists in url
		if exists {
			temp, ok := coll.Collections[col[0]]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("collection not found")})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"schema": temp.Schema})
			return
		}
		schemas := make(map[string]*config.TableRule)
		for p, val := range coll.Collections {
			schemas[p] = &config.TableRule{Schema: val.Schema}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"schemas": schemas})
	}
}

// HandleCollectionRules is an endpoint handler which update database collection rules in config & creates collection if it doesn't exist
func HandleCollectionRules(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		if err := syncman.SetCollectionRules(ctx, project, dbType, col, &v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}

//HandleGetCollectionRules returns handler to get collection rule
func HandleGetCollectionRules(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id and dbType
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbType, exists := r.URL.Query()["dbType"]

		if !exists {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("dbType not provided in url")})
		}

		//gel collection id
		col, exists := r.URL.Query()["col"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get databaseConfig
		databaseConfig, ok := project.Modules.Crud[dbType[0]]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("specified database not present in config")})
			return
		}

		//check col present in url
		if exists {
			//get collection
			collection, ok := databaseConfig.Collections[col[0]]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("specified collection not present in config")})
				return
			}
			collectionsRule := make(map[string]*config.TableRule)
			collectionsRule[col[0]] = &config.TableRule{IsRealTimeEnabled: collection.IsRealTimeEnabled, Rules: collection.Rules}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"rule": collectionsRule})
			return
		}

		collectionsRules := make(map[string]*config.TableRule)
		for p, val := range databaseConfig.Collections {
			collectionsRules[p] = &config.TableRule{IsRealTimeEnabled: val.IsRealTimeEnabled, Rules: val.Rules}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"rules": collectionsRules})
	}
}

// HandleReloadSchema is an endpoint handler which return & sets the schemas of all collection in config
func HandleReloadSchema(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		colResult, err := syncman.SetReloadSchema(ctx, dbType, project, schemaArg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"collections": colResult})
		// return
	}
}

// HandleCreateProject is an endpoint handler which adds a project configuration in config
func HandleCreateProject(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		projectConfig := config.Project{}
		_ = json.NewDecoder(r.Body).Decode(&projectConfig)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		statusCode, err := syncman.CreateProjectConfig(ctx, &projectConfig)
		if err != nil {
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}

// HandleSchemaInspection gets the schema for particular collection & update the database collection schema in config
func HandleSchemaInspection(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		col := vars["col"]
		project := vars["project"]

		schema, err := schemaArg.SchemaInspection(ctx, dbType, project, col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetSchemaInspection(ctx, project, dbType, col, schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"schema": schema})
		// return
	}
}

// HandleModifyAllSchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifyAllSchema(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.CrudStub{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncman.SetModifyAllSchema(ctx, dbType, project, schemaArg, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"statue": true})
		// return
	}
}
