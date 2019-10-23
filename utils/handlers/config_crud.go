package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth/schema"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
	"log"
	"net/http"
	"strings"
	"time"
)

// HandleGetCollections is an endpoint handler which return all the collection name for specified data base
func HandleGetCollections(adminMan *admin.Manager, crud *crud.Module, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		project := vars["project"]

		conf, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		collectionsMap := map[string][]string{}

		for dbType, stub := range conf.Modules.Crud {
			if stub.Enabled {
				collections, err := crud.GetCollections(ctx, project, dbType)
				if err != nil {
					log.Println("Get collections error:", err)
					continue
				}

				cols := make([]string, len(collections))
				for i, value := range collections {
					cols[i] = value.TableName
				}

				collectionsMap[dbType] = cols
			}
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(collectionsMap)
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleDeleteCollection(adminMan *admin.Manager, crud *crud.Module, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		projectConfig, err := syncman.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		coll, ok := projectConfig.Modules.Crud[dbType]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": errors.New("error could not find " + dbType + " in crud").Error()})
			return
		}
		delete(coll.Collections, col)
		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleDatabaseConnection(adminMan *admin.Manager, crud *crud.Module, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		type databaseConnection struct {
			connection string `json:"connection"`
			enabled    bool   `json:"enable"`
		}
		v := &databaseConnection{}
		json.NewDecoder(r.Body).Decode(v)

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := crud.Connect(ctx, dbType, v.connection, v.enabled); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		projectConfig, err := syncman.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		coll, ok := projectConfig.Modules.Crud[dbType]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": errors.New("error could not find " + dbType + " in crud").Error()})
			return
		}
		coll.Conn = v.connection
		coll.Enabled = v.enabled
		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"statue": true})
		return
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleEnforceSchema(adminMan *admin.Manager, crud *crud.Module, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		type schema struct {
			schema string `json:"schema"`
		}
		v := &schema{}
		json.NewDecoder(r.Body).Decode(v)

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		//// Create a context of execution
		//ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		//defer cancel()

		projectConfig, err := syncman.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		collection, ok := projectConfig.Modules.Crud[dbType]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": errors.New("error could not find " + dbType + " in crud").Error()})
			return
		}
		collection.Collections[col].Schema = v.schema
		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"statue": true})
		return
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleCollectionRules(adminMan *admin.Manager, crud *crud.Module, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		type collectionConfig struct {
			IsRealTimeEnabled bool                    `json:"isRealtimeEnabled"`
			Rules             map[string]*config.Rule `json:"rules"`
		}
		v := &collectionConfig{}
		json.NewDecoder(r.Body).Decode(v)

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		projectConfig, err := syncman.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		collection, ok := projectConfig.Modules.Crud[dbType]
		if !ok {
			collection.Collections[col] = &config.TableRule{Schema: "", IsRealTimeEnabled: v.IsRealTimeEnabled, Rules: v.Rules}
		} else {
			collection.Collections[col].IsRealTimeEnabled = v.IsRealTimeEnabled
			collection.Collections[col].Rules = v.Rules
		}

		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleReloadSchema(adminMan *admin.Manager, crud *crud.Module, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]
		col := vars["col"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		projectConfig, err := syncman.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		collectionConfig, ok := projectConfig.Modules.Crud[dbType]
		if !ok {

		}
		for _, colValue := range collectionConfig.Collections {
			result, err := schemaArg.SchemaInspection(ctx, dbType, project, col)
			if err != nil {

			}
			colValue.Schema = result
		}

		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database
func HandleCreateProject(adminMan *admin.Manager, crud *crud.Module, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		projectConfig := &config.Project{}
		json.NewDecoder(r.Body).Decode(projectConfig)

		if err := syncman.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}
