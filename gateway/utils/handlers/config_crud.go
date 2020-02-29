package handlers

import (
	"context"
	"encoding/json"
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
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		collections, err := crud.GetCollections(ctx, projectID, dbAlias)
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
		dbAlias := vars["dbAlias"]

		connState := crud.GetConnectionState(ctx, dbAlias)

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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		if err := crud.DeleteTable(ctx, projectID, dbAlias, col); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetDeleteCollection(ctx, projectID, dbAlias, col); err != nil {
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		if err := syncman.SetDatabaseConnection(ctx, projectID, dbAlias, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		if err := syncman.RemoveDatabaseConfig(ctx, projectID, dbAlias); err != nil {
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if err := schemaArg.SchemaModifyAll(ctx, dbAlias, projectID, map[string]*config.TableRule{col: &v}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetModifySchema(ctx, projectID, dbAlias, col, v.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		if err := syncman.SetCollectionRules(ctx, projectID, dbAlias, col, &v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		colResult, err := syncman.SetReloadSchema(ctx, dbAlias, projectID, schemaArg)
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

// HandleInspectCollectionSchema gets the schema for particular collection & update the database collection schema in config
func HandleInspectCollectionSchema(adminMan *admin.Manager, schemaArg *schema.Schema, syncman *syncman.Manager) http.HandlerFunc {
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
		dbAlias := vars["dbAlias"]
		col := vars["col"]
		projectID := vars["project"]

		schema, err := schemaArg.SchemaInspection(ctx, dbAlias, projectID, col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetSchemaInspection(ctx, projectID, dbAlias, col, schema); err != nil {
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
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncman.SetModifyAllSchema(ctx, dbAlias, projectID, schemaArg, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"statue": true})
		// return
	}
}
