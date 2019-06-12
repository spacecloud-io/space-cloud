package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/utils"
)

type requestMetaData struct {
	project, dbType, col, token string
}

// HandleCrudCreate creates the create operation endpoint
func HandleCrudCreate(isProd bool, auth *auth.Module, crud *crud.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.CreateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Check if the user is authenticated
		status, err := auth.IsCreateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the write operation
		err = crud.Create(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message in dev mode
		realtime.SendCreate(meta.dbType, meta.col, &req)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleCrudRead creates the read operation endpoint
func HandleCrudRead(auth *auth.Module, crud *crud.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.ReadRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create empty read options if it does not exist
		if req.Options == nil {
			req.Options = new(model.ReadOptions)
		}

		// Check if the user is authenticated
		status, err := auth.IsReadOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the read operation
		result, err := crud.Read(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	}
}

// HandleCrudUpdate creates the update operation endpoint
func HandleCrudUpdate(isProd bool, auth *auth.Module, crud *crud.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.UpdateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, err := auth.IsUpdateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the update operation
		err = crud.Update(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message in dev mode
		realtime.SendUpdate(ctx, meta.project, meta.dbType, meta.col, &req, crud)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleCrudDelete creates the delete operation endpoint
func HandleCrudDelete(isProd bool, auth *auth.Module, crud *crud.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.DeleteRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, err := auth.IsDeleteOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the delete operation
		err = crud.Delete(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message in dev mode
		realtime.SendDelete(meta.dbType, meta.col, &req)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleCrudAggregate creates the aggregate operation endpoint
func HandleCrudAggregate(auth *auth.Module, crud *crud.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.AggregateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, err := auth.IsAggregateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the aggregate operation
		result, err := crud.Aggregate(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	}
}

func getRequestMetaData(r *http.Request) *requestMetaData {
	// Get the path parameters
	vars := mux.Vars(r)
	project := vars["project"]
	dbType := vars["dbType"]
	col := vars["col"]

	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	token := strings.TrimPrefix(tokens[0], "Bearer ")

	return &requestMetaData{project: project, dbType: dbType, col: col, token: token}
}

// HandleCrudBatch creates the batch operation endpoint
func HandleCrudBatch(isProd bool, auth *auth.Module, crud *crud.Module, realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		var txRequest model.BatchRequest
		json.NewDecoder(r.Body).Decode(&txRequest)
		defer r.Body.Close()

		for _, req := range txRequest.Requests {

			// Make status and error variables
			var status int
			var err error

			switch req.Type {
			case string(utils.Create):
				req := model.CreateRequest{Document: req.Document, Operation: req.Operation}
				status, err = auth.IsCreateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)

			case string(utils.Update):
				req := model.UpdateRequest{Find: req.Find, Update: req.Update, Operation: req.Operation}
				status, err = auth.IsUpdateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)

			case string(utils.Delete):
				req := model.DeleteRequest{Find: req.Find, Operation: req.Operation}
				status, err = auth.IsDeleteOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
			}

			// Send error response
			if err != nil {
				w.WriteHeader(status)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}

		// Perform the batch operation
		err := crud.Batch(ctx, meta.dbType, meta.project, &txRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if !isProd {
			for _, req := range txRequest.Requests {
				switch req.Type {
				case string(utils.Create):
					// Send realtime message in dev mode
					r := model.CreateRequest{Document: req.Document, Operation: req.Operation}
					realtime.SendCreate(meta.dbType, meta.col, &r)

				case string(utils.Update):
					// Send realtime message in dev mode
					r := model.UpdateRequest{Operation: req.Operation, Find: req.Find, Update: req.Update}
					realtime.SendUpdate(ctx, meta.project, meta.dbType, meta.col, &r, crud)

				case string(utils.Delete):
					// Send realtime message in dev mode
					r := model.DeleteRequest{Find: req.Find, Operation: req.Operation}
					realtime.SendDelete(meta.dbType, meta.col, &r)
				}
			}
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
