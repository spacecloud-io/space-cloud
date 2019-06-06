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
		if !isProd {
			var rows []interface{}
			switch req.Operation {
			case utils.One:
				rows = []interface{}{req.Document}
			case utils.All:
				rows = req.Document.([]interface{})
			default:
				rows = []interface{}{}
			}

			for _, t := range rows {
				data := t.(map[string]interface{})

				idVar := "id"
				if meta.dbType == string(utils.Mongo) {
					idVar = "_id"
				}

				// Send realtime message if id fields exists
				if idTemp, p := data[idVar]; p {
					if id, ok := idTemp.(string); ok {
						realtime.Send(&model.FeedData{
							Group:     meta.col,
							DBType:    meta.dbType,
							Type:      utils.RealtimeWrite,
							TimeStamp: time.Now().Unix(),
							DocID:     id,
							Payload:   data,
						})
					}
				}
			}
		}

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
		if !isProd && req.Operation == utils.One {

			idVar := "id"
			if meta.dbType == string(utils.Mongo) {
				idVar = "_id"
			}

			if idTemp, p := req.Find[idVar]; p {
				if id, ok := idTemp.(string); ok {
					// Create the find object
					find := map[string]interface{}{idVar: id}

					data, err := crud.Read(ctx, meta.dbType, meta.project, meta.col, &model.ReadRequest{Find: find, Operation: utils.One})
					if err == nil {
						realtime.Send(&model.FeedData{
							Group:     meta.col,
							Type:      utils.RealtimeWrite,
							TimeStamp: time.Now().Unix(),
							DocID:     id,
							DBType:    meta.dbType,
							Payload:   data.(map[string]interface{}),
						})
					}
				}
			}
		}

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
		if !isProd && req.Operation == utils.One {
			idVar := "id"
			if meta.dbType == string(utils.Mongo) {
				idVar = "_id"
			}

			if idTemp, p := req.Find[idVar]; p {
				if id, ok := idTemp.(string); ok {
					realtime.Send(&model.FeedData{
						Group:     meta.col,
						Type:      utils.RealtimeDelete,
						TimeStamp: time.Now().Unix(),
						DocID:     id,
						DBType:    meta.dbType,
					})
				}
			}
		}

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
					var rows []interface{}
					switch req.Operation {
					case utils.One:
						rows = []interface{}{req.Document}
					case utils.All:
						rows = req.Document.([]interface{})
					default:
						rows = []interface{}{}
					}

					for _, t := range rows {
						data := t.(map[string]interface{})

						idVar := "id"
						if meta.dbType == string(utils.Mongo) {
							idVar = "_id"
						}

						// Send realtime message if id fields exists
						if id, p := data[idVar]; p {
							realtime.Send(&model.FeedData{
								Group:     req.Col,
								DBType:    meta.dbType,
								Type:      utils.RealtimeWrite,
								TimeStamp: time.Now().Unix(),
								DocID:     id.(string),
								Payload:   data,
							})
						}
					}

				case string(utils.Delete):
					if req.Operation == utils.One {
						idVar := "id"
						if meta.dbType == string(utils.Mongo) {
							idVar = "_id"
						}

						if id, p := req.Find[idVar]; p {
							if err != nil {
								realtime.Send(&model.FeedData{
									Group:     req.Col,
									Type:      utils.RealtimeDelete,
									TimeStamp: time.Now().Unix(),
									DocID:     id.(string),
									DBType:    meta.dbType,
								})
							}
						}
					}

				case string(utils.Update):
					// Send realtime message in dev mode
					if req.Operation == utils.One {

						idVar := "id"
						if meta.dbType == string(utils.Mongo) {
							idVar = "_id"
						}

						if id, p := req.Find[idVar]; p {
							// Create the find object
							find := map[string]interface{}{idVar: id}

							data, err := crud.Read(ctx, meta.dbType, meta.project, req.Col, &model.ReadRequest{Find: find, Operation: utils.One})
							if err == nil {
								realtime.Send(&model.FeedData{
									Group:     req.Col,
									Type:      utils.RealtimeWrite,
									TimeStamp: time.Now().Unix(),
									DocID:     id.(string),
									DBType:    meta.dbType,
									Payload:   data.(map[string]interface{}),
								})
							}
						}
					}
				}
			}
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
