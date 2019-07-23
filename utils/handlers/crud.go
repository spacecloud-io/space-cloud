package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

type requestMetaData struct {
	project, dbType, col, token string
}

// HandleCrudCreate creates the create operation endpoint
func HandleCrudCreate(isProd bool, projects *projects.Projects) http.HandlerFunc {
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

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		status, err := state.Auth.IsCreateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message intent
		msgID := state.Realtime.SendCreateIntent(meta.project, meta.dbType, meta.col, &req)

		// Perform the write operation
		err = state.Crud.Create(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			// Send realtime nack
			state.Realtime.SendAck(msgID, meta.project, meta.col, false)

			// Send http response
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message ack
		state.Realtime.SendAck(msgID, meta.project, meta.col, true)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
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
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create empty read options if it does not exist
		if req.Options == nil {
			req.Options = new(model.ReadOptions)
		}

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		status, err := state.Auth.IsReadOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the read operation
		result, err := state.Crud.Read(ctx, meta.dbType, meta.project, meta.col, &req)
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
func HandleCrudUpdate(isProd bool, projects *projects.Projects) http.HandlerFunc {
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

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsUpdateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message intent
		msgID := state.Realtime.SendUpdateIntent(meta.project, meta.dbType, meta.col, &req)

		// Perform the update operation
		err = state.Crud.Update(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			// Send realtime nack
			state.Realtime.SendAck(msgID, meta.project, meta.col, false)

			// Send http response
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message ack
		state.Realtime.SendAck(msgID, meta.project, meta.col, true)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleCrudDelete creates the delete operation endpoint
func HandleCrudDelete(isProd bool, projects *projects.Projects) http.HandlerFunc {
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

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsDeleteOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message intent
		msgID := state.Realtime.SendDeleteIntent(meta.project, meta.dbType, meta.col, &req)

		// Perform the delete operation
		err = state.Crud.Delete(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			// Send realtime nack
			state.Realtime.SendAck(msgID, meta.project, meta.col, false)

			// Send http response
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime message ack
		state.Realtime.SendAck(msgID, meta.project, meta.col, true)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
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
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status, err := state.Auth.IsAggregateOpAuthorised(meta.project, meta.dbType, meta.col, meta.token, &req)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the read operation
		result, err := state.Crud.Aggregate(ctx, meta.dbType, meta.project, meta.col, &req)
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
func HandleCrudBatch(isProd bool, projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		var txRequest model.BatchRequest
		json.NewDecoder(r.Body).Decode(&txRequest)
		defer r.Body.Close()

		type msg struct {
			id, col string
		}

		msgIDs := make([]*msg, len(txRequest.Requests))

		for i, req := range txRequest.Requests {

			// Make status and error variables
			var status int
			var err error

			switch req.Type {
			case string(utils.Create):
				r := model.CreateRequest{Document: req.Document, Operation: req.Operation}
				status, err = state.Auth.IsCreateOpAuthorised(meta.project, meta.dbType, req.Col, meta.token, &r)
				if err == nil {
					// Send realtime message intent if user is authorised (no error)
					msgID := state.Realtime.SendCreateIntent(meta.project, meta.dbType, req.Col, &r)
					msgIDs[i] = &msg{id: msgID, col: req.Col}
				}

			case string(utils.Update):
				r := model.UpdateRequest{Find: req.Find, Update: req.Update, Operation: req.Operation}
				status, err = state.Auth.IsUpdateOpAuthorised(meta.project, meta.dbType, req.Col, meta.token, &r)
				if err == nil {
					// Send realtime message intent if user is authorised (no error)
					msgID := state.Realtime.SendUpdateIntent(meta.project, meta.dbType, req.Col, &r)
					msgIDs[i] = &msg{id: msgID, col: req.Col}
				}

			case string(utils.Delete):
				r := model.DeleteRequest{Find: req.Find, Operation: req.Operation}
				status, err = state.Auth.IsDeleteOpAuthorised(meta.project, meta.dbType, req.Col, meta.token, &r)
				if err == nil {
					// Send realtime message intent if user is authorised (no error)
					msgID := state.Realtime.SendDeleteIntent(meta.project, meta.dbType, req.Col, &r)
					msgIDs[i] = &msg{id: msgID, col: req.Col}
				}
			}

			// Send error response
			if err != nil {
				// Send realtime nack
				for j := 0; j < i; j++ {
					state.Realtime.SendAck(msgIDs[j].id, meta.project, msgIDs[j].col, false)
				}

				// Send http response
				w.WriteHeader(status)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}

		// Perform the batch operation
		err = state.Crud.Batch(ctx, meta.dbType, meta.project, &txRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Send realtime nack
		for _, m := range msgIDs {
			state.Realtime.SendAck(m.id, meta.project, m.col, false)
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
