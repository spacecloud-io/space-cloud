package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleGetCollections is an endpoint handler which return all the collection(table) names for specified data base
func HandleGetCollections(adminMan *admin.Manager, projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer r.Body.Close()

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
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		collections, err := state.Crud.GetCollections(ctx, projectID, dbAlias)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		cols := make([]string, len(collections))
		for i, value := range collections {
			cols[i] = value.TableName
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"collections": cols})
	}
}

// HandleGetConnectionState gives the status of connection state of client
func HandleGetConnectionState(adminMan *admin.Manager, projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		defer r.Body.Close()

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		connState := state.Crud.GetConnectionState(ctx, dbAlias)

		w.WriteHeader(http.StatusOK) // http status code
		_ = json.NewEncoder(w).Encode(map[string]bool{"status": connState})
	}
}

// HandleDeleteCollection is an endpoint handler which deletes a table in specified database & removes it from config
func HandleDeleteCollection(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := state.Crud.DeleteTable(ctx, projectID, dbAlias, col); err != nil {
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
func HandleDatabaseConnection(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
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
func HandleRemoveDatabaseConfig(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		if err := syncman.RemoveDatabaseConfig(ctx, projectID, dbAlias); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleModifySchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifySchema(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		json.NewDecoder(r.Body).Decode(&v)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if err := state.Schema.SchemaModifyAll(ctx, dbAlias, projectID, map[string]*config.TableRule{col: &v}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetModifySchema(ctx, projectID, dbAlias, col, v.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleCollectionRules is an endpoint handler which update database collection rules in config & creates collection if it doesn't exist
func HandleCollectionRules(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		json.NewDecoder(r.Body).Decode(&v)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleReloadSchema is an endpoint handler which return & sets the schemas of all collection in config
func HandleReloadSchema(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		colResult, err := syncman.SetReloadSchema(ctx, dbAlias, projectID, state.Schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"collections": colResult})
		return
	}
}

// HandleCreateProject is an endpoint handler which adds a project configuration in config
func HandleCreateProject(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		projectConfig := config.Project{}
		json.NewDecoder(r.Body).Decode(&projectConfig)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err, statusCode := syncman.CreateProjectConfig(ctx, &projectConfig)
		if err != nil {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

// HandleSchemaInspection gets the schema for particular collection & update the database collection schema in config
func HandleSchemaInspection(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

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
		dbAlias := vars["dbAlias"]
		col := vars["col"]
		projectID := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		schema, err := state.Schema.SchemaInspection(ctx, dbAlias, projectID, col)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncman.SetSchemaInspection(ctx, projectID, dbAlias, col, schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"schema": schema})
		return
	}
}

// HandleModifySchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifyAllSchema(adminMan *admin.Manager, projects *projects.Projects, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.CrudStub{}
		json.NewDecoder(r.Body).Decode(&v)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncman.SetModifyAllSchema(ctx, dbAlias, projectID, state.Schema, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{"statue": true})
		return
	}
}
