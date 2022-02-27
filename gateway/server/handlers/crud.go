package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	authHelpers "github.com/spaceuptech/space-cloud/gateway/modules/auth/helpers"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type requestMetaData struct {
	projectID, dbType, col, token string
}

// HandleCrudPreparedQuery creates the PreparedQuery operation endpoint
func HandleCrudPreparedQuery(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		project := vars["project"]
		id := vars["id"]
		token := utils.GetTokenFromHeader(r)

		auth, err := modules.Auth(project)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		crud, err := modules.DB(project)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		// Load the request from the body
		req := model.PreparedQueryRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Check if the user is authenticated
		actions, reqParams, err := auth.IsPreparedQueryAuthorised(ctx, project, dbAlias, id, token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the PreparedQuery operation
		result, _, err := crud.ExecPreparedQuery(ctx, dbAlias, id, &req, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// function to do postProcessing on result
		_ = authHelpers.PostProcessMethod(ctx, auth.GetAESKey(), actions, result)

		// Give positive acknowledgement
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": result})
	}
}

// HandleCrudCreate creates the create operation endpoint
func HandleCrudCreate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		meta := getRequestMetaData(r)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.CreateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Check if the user is authenticated
		reqParams, err := auth.IsCreateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the write operation
		err = crud.Create(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {

			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudRead creates the read operation endpoint
func HandleCrudRead(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.ReadRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Create empty read options if it does not exist
		if req.Options == nil {
			req.Options = new(model.ReadOptions)
		}

		// Check if read op is authorised
		// NOTE: meta.dbType is actually the dbAlias
		dbType, _ := crud.GetDBType(meta.dbType)

		returnWhere := model.ReturnWhereStub{Col: meta.col, PrefixColName: len(req.Options.Join) > 0, ReturnWhere: dbType != string(model.Mongo), Where: map[string]interface{}{}}
		actions, reqParams, err := auth.IsReadOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req, returnWhere)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}
		if len(returnWhere.Where) > 0 {
			req.MatchWhere = append(req.MatchWhere, returnWhere.Where)
		}

		if req.PostProcess == nil {
			req.PostProcess = map[string]*model.PostProcess{}
		}
		req.PostProcess[meta.col] = actions

		if err := auth.RunAuthForJoins(ctx, meta.projectID, dbType, meta.dbType, meta.token, &req, req.Options.Join); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		result, _, err := crud.Read(ctx, meta.dbType, meta.col, &req, reqParams)
		// Perform the read operation

		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// function to do postProcessing on result
		_ = authHelpers.PostProcessMethod(ctx, auth.GetAESKey(), actions, result)

		// Give positive acknowledgement
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": result})
	}
}

// HandleCrudUpdate creates the update operation endpoint
func HandleCrudUpdate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		meta := getRequestMetaData(r)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.UpdateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := auth.IsUpdateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the update operation
		err = crud.Update(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {

			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudDelete creates the delete operation endpoint
func HandleCrudDelete(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		meta := getRequestMetaData(r)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.DeleteRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := auth.IsDeleteOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the delete operation
		err = crud.Delete(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {
			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudAggregate creates the aggregate operation endpoint
func HandleCrudAggregate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		meta := getRequestMetaData(r)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.AggregateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := auth.IsAggregateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the aggregate operation
		result, err := crud.Aggregate(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": result})
	}
}

func getRequestMetaData(r *http.Request) *requestMetaData {
	// Get the path parameters
	vars := mux.Vars(r)
	projectID := vars["project"]
	dbAlias := vars["dbAlias"]
	col := vars["col"]

	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	token := strings.TrimPrefix(tokens[0], "Bearer ")

	return &requestMetaData{projectID: projectID, dbType: dbAlias, col: col, token: token}
}

// HandleCrudBatch creates the batch operation endpoint
func HandleCrudBatch(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		meta := getRequestMetaData(r)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		auth, err := modules.Auth(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		crud, err := modules.DB(meta.projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		var txRequest model.BatchRequest
		_ = json.NewDecoder(r.Body).Decode(&txRequest)
		defer utils.CloseTheCloser(r.Body)

		var reqParams model.RequestParams
		for _, req := range txRequest.Requests {
			// Make error variables
			var err error

			switch req.Type {
			case string(model.Create):
				r := model.CreateRequest{Document: req.Document, Operation: req.Operation}
				reqParams, err = auth.IsCreateOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			case string(model.Update):
				r := model.UpdateRequest{Find: req.Find, Update: req.Update, Operation: req.Operation}
				reqParams, err = auth.IsUpdateOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			case string(model.Delete):
				r := model.DeleteRequest{Find: req.Find, Operation: req.Operation}
				reqParams, err = auth.IsDeleteOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			}

			// Send error response
			if err != nil {
				// Send http response
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
				return
			}
		}

		reqParams.Resource = "db-batch"
		reqParams = utils.ExtractRequestParams(r, reqParams, txRequest)

		err = crud.Batch(ctx, meta.dbType, &txRequest, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
