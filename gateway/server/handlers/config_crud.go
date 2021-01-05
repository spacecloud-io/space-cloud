package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleGetAllTableNames is an endpoint handler which return all the collection(table) names for specified data base
func HandleGetAllTableNames(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "db-config", "read", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		crud, err := modules.DB(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		collections, err := crud.GetCollections(ctx, dbAlias)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		cols := make([]string, len(collections))
		for i, value := range collections {
			cols[i] = value.TableName
		}
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: cols})
	}
}

// HandleGetDatabaseConnectionState gives the status of connection state of client
func HandleGetDatabaseConnectionState(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "db-config", "read", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		crud, err := modules.DB(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		connState := crud.GetConnectionState(ctx, dbAlias)

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: connState})
	}
}

// HandleDeleteTable is an endpoint handler which deletes a table in specified database & removes it from config
func HandleDeleteTable(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		crud, err := modules.DB(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, err := syncman.SetDeleteCollection(ctx, projectID, dbAlias, col, crud, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleSetDatabaseConfig is an endpoint handler which updates database config & connects to database
func HandleSetDatabaseConfig(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		v := config.DatabaseConfig{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-config", "modify", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, v)
		status, err := syncman.SetDatabaseConnection(ctx, projectID, dbAlias, &v, reqParams)
		fmt.Println("ABCDEF", err == nil)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetDatabaseConfig returns handler to get Database Collection
func HandleGetDatabaseConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := "*"
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-config", "read", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, dbConfig, err := syncMan.GetDatabaseConfig(ctx, projectID, dbAlias, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: dbConfig})
	}
}

// HandleRemoveDatabaseConfig is an endpoint handler which removes database config
func HandleRemoveDatabaseConfig(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-config", "modify", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, err := syncman.RemoveDatabaseConfig(ctx, projectID, dbAlias, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetPreparedQuery returns handler to get PreparedQuery
func HandleGetPreparedQuery(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := "*"
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		idQuery, exists := r.URL.Query()["id"]
		id := "*"
		if exists {
			id = idQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-prepared-query", "read", map[string]string{"project": projectID, "db": dbAlias, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, result, err := syncMan.GetPreparedQuery(ctx, projectID, dbAlias, id, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: result})
	}
}

// HandleSetPreparedQueries is an endpoint handler which updates database PreparedQueries
func HandleSetPreparedQueries(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		id := vars["id"]

		v := config.DatbasePreparedQuery{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-prepared-query", "modify", map[string]string{"project": projectID, "db": dbAlias, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, &v)
		status, err := syncman.SetPreparedQueries(ctx, projectID, dbAlias, id, &v, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleRemovePreparedQueries is an endpoint handler which removes database PreparedQueries
func HandleRemovePreparedQueries(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		id := vars["id"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-prepared-query", "modify", map[string]string{"project": projectID, "db": dbAlias, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, err := syncman.RemovePreparedQueries(ctx, projectID, dbAlias, id, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleModifySchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifySchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		v := config.DatabaseSchema{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, v)
		status, err := syncman.SetModifySchema(ctx, projectID, dbAlias, col, &v, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetSchemas returns handler to get schema
func HandleGetSchemas(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and dbType from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := "*"
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		colQuery, exists := r.URL.Query()["col"]
		col := "*"
		if exists {
			col = colQuery[0]
		}
		formatQuery, exists := r.URL.Query()["format"]
		format := "graphql"
		if exists {
			format = formatQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "read", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, schemas, err := syncMan.GetSchemas(ctx, projectID, dbAlias, col, format, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: schemas})
	}
}

// HandleSetTableRules is an endpoint handler which update database collection rules in config & creates collection if it doesn't exist
func HandleSetTableRules(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		v := config.DatabaseRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-rule", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, v)
		status, err := syncman.SetCollectionRules(ctx, projectID, dbAlias, col, &v, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetTableRules returns handler to get collection rule
func HandleGetTableRules(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and dbAlias
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := "*"
		dbAliasQuery, exists := r.URL.Query()["dbAlias"]
		if exists {
			dbAlias = dbAliasQuery[0]
		}
		col := "*"
		colQuery, exists := r.URL.Query()["col"]
		if exists {
			col = colQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-rule", "read", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, dbConfig, err := syncMan.GetCollectionRules(ctx, projectID, dbAlias, col, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: dbConfig})
	}
}

// HandleDeleteTableRules is an endpoint handler which deletes database collection rules in config
func HandleDeleteTableRules(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]
		col := vars["col"]

		v := config.TableRule{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-rule", "delete", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, v)
		status, err := syncman.DeleteCollectionRules(ctx, projectID, dbAlias, col, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleReloadSchema is an endpoint handler which return & sets the schemas of all collection in config
func HandleReloadSchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": "*"})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, err := syncman.SetReloadSchema(ctx, dbAlias, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleInspectCollectionSchema gets the schema for particular collection & update the database collection schema in config
func HandleInspectCollectionSchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		col := vars["col"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": col})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		logicalDBName, err := syncman.GetLogicalDatabaseName(ctx, projectID, dbAlias)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		schema, err := modules.Schema(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		s, err := schema.SchemaInspection(ctx, dbAlias, logicalDBName, col)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, err := syncman.SetSchemaInspection(ctx, projectID, dbAlias, col, s, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleUntrackCollectionSchema removes the collection from the database config
func HandleUntrackCollectionSchema(adminMan *admin.Manager, modules *modules.Modules, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		col := vars["col"]
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, err := syncman.RemoveCollection(ctx, projectID, dbAlias, col, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleModifyAllSchema is an endpoint handler which updates the existing schema & updates the config
func HandleModifyAllSchema(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		v := config.CrudStub{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "db-schema", "modify", map[string]string{"project": projectID, "db": dbAlias, "col": "*"})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, v)
		status, err := syncman.SetModifyAllSchema(ctx, dbAlias, projectID, v, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleInspectTrackedCollectionsSchema is an endpoint handler which return schema for all tracked collections of a particular database
func HandleInspectTrackedCollectionsSchema(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Create a context of execution
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "db-schema", "read", map[string]string{"project": projectID, "db": dbAlias, "col": "*"})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		schema, err := modules.Schema(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		schemas, err := schema.GetCollectionSchema(ctx, projectID, dbAlias)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: schemas})
	}
}
