package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

type requestMetaData struct {
	project, dbType, col, token string
}

func (s *server) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Check if the user is authicated
		authObj, err := s.auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Create)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
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
		err = s.auth.IsAuthorized(meta.dbType, meta.col, utils.Create, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		// Perform the write operation
		err = s.crud.Create(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func (s *server) handleRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Check if the user is authicated
		authObj, err := s.auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Read)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
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
			"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
			"project": meta.project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = s.auth.IsAuthorized(meta.dbType, meta.col, utils.Read, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		// Perform the read operation
		result, err := s.crud.Read(ctx, meta.dbType, meta.project, meta.col, &req)
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

func (s *server) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Check if the user is authicated
		authObj, err := s.auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Update)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
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
		err = s.auth.IsAuthorized(meta.dbType, meta.col, utils.Update, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		// Perform the update operation
		err = s.crud.Update(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Check if the user is authicated
		authObj, err := s.auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Delete)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
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
		err = s.auth.IsAuthorized(meta.dbType, meta.col, utils.Delete, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		// Perform the delete operation
		err = s.crud.Delete(ctx, meta.dbType, meta.project, meta.col, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func (s *server) handleAggregate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		meta := getRequestMetaData(r)

		// Check if the user is authicated
		authObj, err := s.auth.IsAuthenticated(meta.token, meta.dbType, meta.col, utils.Aggregation)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
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
		err = s.auth.IsAuthorized(meta.dbType, meta.col, utils.Aggregation, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		// Perform the read operation
		result, err := s.crud.Aggregate(ctx, meta.dbType, meta.project, meta.col, &req)
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

	return &requestMetaData{project: project, dbType: dbType, col: col, token: tokens[0]}
}
