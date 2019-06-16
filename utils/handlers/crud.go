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

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Create)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.CreateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"doc": &req.Document, "op": req.Operation, "auth": authObj},
			"project": meta.project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = state.Auth.IsAuthorized(meta.project, meta.dbType, meta.col, utils.Create, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the write operation
		err = state.Crud.Create(ctx, meta.dbType, meta.project, meta.col, &req)
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
						state.Realtime.Send(&model.FeedData{
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
func HandleCrudRead(projects *projects.Projects) http.HandlerFunc {
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

		// Check if the user is authenticated
		authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Read)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.ReadRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create empty read options if it does not exist
		if req.Options == nil {
			req.Options = new(model.ReadOptions)
		}

		// Create an args object
		args := map[string]interface{}{
			"args": map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
		}

		// Check if user is authorized to make this request
		err = state.Auth.IsAuthorized(meta.project, meta.dbType, meta.col, utils.Read, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
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

		// Load the project state
		state, err := projects.LoadProject(meta.project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Check if the user is authenticated
		authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Update)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.UpdateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"find": req.Find, "update": req.Update, "op": req.Operation, "auth": authObj},
			"project": meta.project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = state.Auth.IsAuthorized(meta.project, meta.dbType, meta.col, utils.Update, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the update operation
		err = state.Crud.Update(ctx, meta.dbType, meta.project, meta.col, &req)
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

					data, err := state.Crud.Read(ctx, meta.dbType, meta.project, meta.col, &model.ReadRequest{Find: find, Operation: utils.One})
					if err == nil {
						state.Realtime.Send(&model.FeedData{
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
func HandleCrudDelete(isProd bool, projects *projects.Projects) http.HandlerFunc {
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

		// Check if the user is authenticated
		authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Delete)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.DeleteRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
			"project": meta.project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = state.Auth.IsAuthorized(meta.project, meta.dbType, meta.col, utils.Delete, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Perform the delete operation
		err = state.Crud.Delete(ctx, meta.dbType, meta.project, meta.col, &req)
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
					state.Realtime.Send(&model.FeedData{
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
func HandleCrudAggregate(projects *projects.Projects) http.HandlerFunc {
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

		// Check if the user is authicated
		authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Aggregation)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := model.AggregateRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"find": req.Pipeline, "op": req.Operation, "auth": authObj},
			"project": meta.project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = state.Auth.IsAuthorized(meta.project, meta.dbType, meta.col, utils.Aggregation, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
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

		args := map[string]interface{}{}
		for _, req := range txRequest.Requests {

			switch req.Type {
			case string(utils.Update):
				authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, req.Col, utils.Update)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
					return
				}
				args = map[string]interface{}{
					"args":    map[string]interface{}{"find": req.Find, "update": req.Update, "op": req.Operation, "auth": authObj},
					"project": meta.project, // Don't forget to do this for every request
				}

				// Check if user is authorized to make this request
				err = state.Auth.IsAuthorized(meta.project, meta.dbType, req.Col, utils.Update, args)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
					return
				}

			case string(utils.Create):
				authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, req.Col, utils.Create)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
					return
				}
				// Create an args object
				args = map[string]interface{}{
					"args":    map[string]interface{}{"doc": &req.Document, "op": req.Operation, "auth": authObj},
					"project": meta.project, // Don't forget to do this for every request
				}

				// Check if user is authorized to make this request
				err = state.Auth.IsAuthorized(meta.project, meta.dbType, req.Col, utils.Create, args)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
					return
				}

			case string(utils.Delete):

				authObj, err := state.Auth.IsAuthenticated(meta.token, meta.dbType, req.Col, utils.Delete)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
					return
				}
				// Create an args object
				args = map[string]interface{}{
					"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
					"project": meta.project, // Don't forget to do this for every request
				}

				// Check if user is authorized to make this request
				err = state.Auth.IsAuthorized(meta.project, meta.dbType, req.Col, utils.Delete, args)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
					return
				}

			}
		}

		// Perform the batch operation
		err = state.Crud.Batch(ctx, meta.dbType, meta.project, &txRequest)
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
							state.Realtime.Send(&model.FeedData{
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
								state.Realtime.Send(&model.FeedData{
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

							data, err := state.Crud.Read(ctx, meta.dbType, meta.project, req.Col, &model.ReadRequest{Find: find, Operation: utils.One})
							if err == nil {
								state.Realtime.Send(&model.FeedData{
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
