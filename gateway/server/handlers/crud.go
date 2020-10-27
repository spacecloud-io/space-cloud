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
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type requestMetaData struct {
	projectID, dbType, col, token string
}

// HandleCrudPreparedQuery creates the PreparedQuery operation endpoint
func HandleCrudPreparedQuery(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		project := vars["project"]
		id := vars["id"]
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := model.PreparedQueryRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Check if the user is authenticated
		actions, reqParams, err := auth.IsPreparedQueryAuthorised(ctx, project, dbAlias, id, token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the PreparedQuery operation
		result, err := crud.ExecPreparedQuery(ctx, dbAlias, id, &req, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// function to do postProcessing on result
		_ = auth.PostProcessMethod(ctx, actions, result)

		// Give positive acknowledgement
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": result})
	}
}

// HandleCrudCreate creates the create operation endpoint
func HandleCrudCreate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.CreateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the user is authenticated
		reqParams, err := auth.IsCreateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the write operation
		err = crud.Create(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {

			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudRead creates the read operation endpoint
func HandleCrudRead(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.ReadRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Create empty read options if it does not exist
		if req.Options == nil {
			req.Options = new(model.ReadOptions)
		}

		// Rest API is not allowed to do joins for security reasons
		req.Options.Join = nil

		// Check if the user is authenticated
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		actions, reqParams, err := auth.IsReadOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req, model.ReturnWhereStub{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		result, err := crud.Read(ctx, meta.dbType, meta.col, &req, reqParams)
		// Perform the read operation

		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// function to do postProcessing on result
		_ = auth.PostProcessMethod(ctx, actions, result)

		// Give positive acknowledgement
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": result})
	}
}

// HandleCrudUpdate creates the update operation endpoint
func HandleCrudUpdate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.UpdateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := auth.IsUpdateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the update operation
		err = crud.Update(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {

			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudDelete creates the delete operation endpoint
func HandleCrudDelete(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.DeleteRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := auth.IsDeleteOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the delete operation
		err = crud.Delete(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {
			// Send http response
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleCrudAggregate creates the aggregate operation endpoint
func HandleCrudAggregate(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.AggregateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := auth.IsAggregateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		// Perform the aggregate operation
		result, err := crud.Aggregate(ctx, meta.dbType, meta.col, &req, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
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

		auth := modules.Auth()
		crud := modules.DB()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		var txRequest model.BatchRequest
		_ = json.NewDecoder(r.Body).Decode(&txRequest)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

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
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
				return
			}
		}

		reqParams.Resource = "db-batch"
		reqParams = utils.ExtractRequestParams(r, reqParams, txRequest)

		err := crud.Batch(ctx, meta.dbType, &txRequest, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Give positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
