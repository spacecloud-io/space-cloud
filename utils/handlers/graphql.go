package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth/schema"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/graphql"
)

// HandleGraphQLRequest creates the graphql operation endpoint
func HandleGraphQLRequest(graphql *graphql.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		_, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		pid := graphql.GetProjectID()

		if projectID != pid {
			//throw some error
			w.WriteHeader(http.StatusInternalServerError) //http status codee
			json.NewEncoder(w).Encode(map[string]string{"error": "project id doesn't match"})
			return
		}

		// Get the path parameters
		token := getRequestMetaData(r).token
		// Load the request from the body
		req := model.GraphQLRequest{}

		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		var wg sync.WaitGroup
		wg.Add(1)

		graphql.ExecGraphQLQuery(&req, token, func(op interface{}, err error) {
			defer wg.Done()

			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}

			w.WriteHeader(http.StatusOK) //http status codee
			json.NewEncoder(w).Encode(map[string]interface{}{"data": op})
			return
		})

		wg.Wait()
	}

}

// HandleInspectionRequest creates the schema inspection endpoint
func HandleInspectionRequest(adminMan *admin.Manager, schemaArg *schema.Schema) http.HandlerFunc {
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

		fmt.Println("schema", schemaArg)

		result, err := schemaArg.SchemaInspection(ctx, dbType, project, col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"schema": result})
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
		w.WriteHeader(http.StatusOK) //http status codee
		return
	}
}

// HandleGetCollections is an endpoint handler which return all the collection name for specified data base
func HandleGetCollections(adminMan *admin.Manager, crud *crud.Module) http.HandlerFunc {
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

		switch utils.DBType(dbType) {
		case utils.Mongo, utils.MySQL, utils.Postgres:
			collections, err := crud.GetCollections(ctx, project, dbType)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			col := make([]string, len(collections))
			for key, value := range collections {
				col[key] = value.TableName
			}
			w.WriteHeader(http.StatusOK) //http status codee
			json.NewEncoder(w).Encode(map[string]interface{}{"collections": col})
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": errors.New("collections wrong database").Error()})
			return
		}
	}
}

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
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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
		json.NewEncoder(w).Encode(map[string]interface{}{"value": schemas})
		return
	}
}
