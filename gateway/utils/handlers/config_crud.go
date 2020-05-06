package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleGetAllTableNames is an endpoint handler which return all the collection(table) names for specified data base
func HandleGetAllTableNames(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		crud := modules.DB()
		collections, err := crud.GetCollections(ctx, projectID, dbAlias)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cols := make([]string, len(collections))
		for i, value := range collections {
			cols[i] = value.TableName
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: cols})
	}
}

// HandleGetDatabaseConnectionState gives the status of connection state of client
func HandleGetDatabaseConnectionState(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]

		crud := modules.DB()
		connState := crud.GetConnectionState(ctx, dbAlias)

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: connState})
	}
}

// HandleDeleteTable is an endpoint handler which deletes a table in specified database & removes it from config
func HandleDeleteTable(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		crud := modules.DB()
		if err := crud.DeleteTable(ctx, projectID, dbAlias, col); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := syncman.SetDeleteCollection(ctx, projectID, dbAlias, col); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleSetDatabaseConfig is an endpoint handler which updates database config & connects to database
func HandleSetDatabaseConfig(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.CrudStub{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		if err := syncman.SetDatabaseConnection(ctx, projectID, dbAlias, v); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetDatabaseConfig returns handler to get Database Collection
func HandleGetDatabaseConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := ""
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		dbConfig, err := syncMan.GetDatabaseConfig(ctx, projectID, dbAlias)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: dbConfig})
	}
}

// HandleRemoveDatabaseConfig is an endpoint handler which removes database config
func HandleRemoveDatabaseConfig(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		if err := syncman.RemoveDatabaseConfig(ctx, projectID, dbAlias); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
		// return
	}
}

// HandleModifySchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifySchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		schema := modules.Schema()
		if err := schema.SchemaModifyAll(ctx, dbAlias, projectID, map[string]*config.TableRule{col: &v}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := syncman.SetModifySchema(ctx, projectID, dbAlias, col, v.Schema); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetSchemas returns handler to get schema
func HandleGetSchemas(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		// get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := ""
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		colQuery, exists := r.URL.Query()["col"]
		col := ""
		if exists {
			col = colQuery[0]
		}
		schemas, err := syncMan.GetSchemas(ctx, projectID, dbAlias, col)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: schemas})
	}
}

// HandleSetTableRules is an endpoint handler which update database collection rules in config & creates collection if it doesn't exist
func HandleSetTableRules(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.TableRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		if err := syncman.SetCollectionRules(ctx, projectID, dbAlias, col, &v); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
		// return
	}
}

// HandleGetTableRules returns handler to get collection rule
func HandleGetTableRules(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		// get project id and dbAlias
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := ""
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		col := ""
		colQuery, exists := r.URL.Query()["col"]
		if exists {
			col = colQuery[0]
		}

		dbConfig, err := syncMan.GetCollectionRules(ctx, projectID, dbAlias, col)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: dbConfig})
	}
}

// HandleReloadSchema is an endpoint handler which return & sets the schemas of all collection in config
func HandleReloadSchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		schema := modules.Schema()
		colResult, err := syncman.SetReloadSchema(ctx, dbAlias, projectID, schema)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: colResult})
		// return
	}
}

// HandleInspectCollectionSchema gets the schema for particular collection & update the database collection schema in config
func HandleInspectCollectionSchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		col := vars["col"]
		projectID := vars["project"]

		schema := modules.Schema()
		s, err := schema.SchemaInspection(ctx, dbAlias, projectID, col)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := syncman.SetSchemaInspection(ctx, projectID, dbAlias, col, s); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: s})
		// return
	}
}

// HandleModifyAllSchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifyAllSchema(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		v := config.CrudStub{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		if err := syncman.SetModifyAllSchema(ctx, dbAlias, projectID, v); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
		// return
	}
}
