package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/modules/auth/schema"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
	"net/http"
	"strings"
	"time"
)

// HandleGetCollectionSchemas is an endpoint handler which return schema for all the collection in the config.crud
func HandleGetCollectionSchemas(adminMan *admin.Manager, schema *schema.Schema) http.HandlerFunc {
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
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Minute)
		defer cancel()

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]

		schemas, err := schema.GetCollectionSchema(ctx, project, dbType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"collections": schemas})
		return
	}
}

// HandleInspectionRequest creates the schema inspection endpoint
func HandleInspectionRequest(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
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
		col := vars["col"]
		project := vars["project"]

		schema, err := schemaArg.SchemaInspection(ctx, dbType, project, col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// update schema in config
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
		coll.Collections[col].Schema = schema
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

// HandleCreationRequest creates the schema inspection endpoint
func HandleCreationRequest(adminMan *admin.Manager, schema *schema.Schema) http.HandlerFunc {
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

		// Load the request from the body
		type schemaRequest struct {
			Schema string `json:"schema"`
		}

		schemaDecode := schemaRequest{}
		json.NewDecoder(r.Body).Decode(&schemaDecode)
		defer r.Body.Close()

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		col := vars["col"]
		project := vars["project"]

		if err := schema.SchemaCreation(ctx, dbType, col, project, schemaDecode.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}
}
