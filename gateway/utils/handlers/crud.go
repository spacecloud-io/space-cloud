package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

type requestMetaData struct {
	projectID, dbType, col, token string
}

// HandleCrudCreate creates the create operation endpoint
func HandleCrudCreate(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.CreateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		status, err := state.Auth.IsCreateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the write operation
		err = state.Crud.Create(ctx, meta.dbType, meta.projectID, meta.col, &req)
		if err != nil {

			// Send http response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleCrudRead creates the read operation endpoint
func HandleCrudRead(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

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

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		actions, status, err := state.Auth.IsReadOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the read operation
		result, err := state.Crud.Read(ctx, meta.dbType, meta.projectID, meta.col, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// function to do postProcessing on result
		_ = state.Auth.PostProcessMethod(actions, result)

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	}
}

// HandleCrudUpdate creates the update operation endpoint
func HandleCrudUpdate(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.UpdateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsUpdateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the update operation
		err = state.Crud.Update(ctx, meta.dbType, meta.projectID, meta.col, &req)
		if err != nil {

			// Send http response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleCrudDelete creates the delete operation endpoint
func HandleCrudDelete(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.DeleteRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsDeleteOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the delete operation
		err = state.Crud.Delete(ctx, meta.dbType, meta.projectID, meta.col, &req)
		if err != nil {
			// Send http response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleCrudAggregate creates the aggregate operation endpoint
func HandleCrudAggregate(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the request from the body
		req := model.AggregateRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsAggregateOpAuthorised(ctx, meta.projectID, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the read operation
		result, err := state.Crud.Aggregate(ctx, meta.dbType, meta.projectID, meta.col, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
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
func HandleCrudBatch(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the project state
		state, err := projects.LoadProject(meta.projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		var txRequest model.BatchRequest
		_ = json.NewDecoder(r.Body).Decode(&txRequest)
		defer utils.CloseTheCloser(r.Body)

		for _, req := range txRequest.Requests {

			// Make status and error variables
			var status int
			var err error

			switch req.Type {
			case string(utils.Create):
				r := model.CreateRequest{Document: req.Document, Operation: req.Operation}
				status, err = state.Auth.IsCreateOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			case string(utils.Update):
				r := model.UpdateRequest{Find: req.Find, Update: req.Update, Operation: req.Operation}
				status, err = state.Auth.IsUpdateOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			case string(utils.Delete):
				r := model.DeleteRequest{Find: req.Find, Operation: req.Operation}
				status, err = state.Auth.IsDeleteOpAuthorised(ctx, meta.projectID, meta.dbType, req.Col, meta.token, &r)

			}

			// Send error response
			if err != nil {
				// Send http response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}

		// Perform the batch operation
		err = state.Crud.Batch(ctx, meta.dbType, meta.projectID, &txRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
