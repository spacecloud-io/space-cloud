package userman

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"


	"github.com/gorilla/mux"
)

// HandleProfile returns the handler for fetching single user profile
func (m *Module) HandleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]
		id := vars["id"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		status, result, err := m.Profile(ctx, token, dbType, project, id)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"user": result})
	}
}

// HandleProfiles returns the handler for fetching all user profiles
func (m *Module) HandleProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		status, result, err := m.Profiles(ctx, token, dbType, project)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"users": result})
	}
}

// HandleEmailSignIn returns the handler for email sign in
func (m *Module) HandleEmailSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := m.EmailSignIn(ctx, dbType, project, req["email"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignUp returns the handler for email sign up
func (m *Module) HandleEmailSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := m.EmailSignUp(ctx, dbType, project, req["email"].(string), req["name"].(string), req["pass"].(string), req["role"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignUp returns the handler for email sign up
func (m *Module) HandleEmailEditProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]
		id := vars["id"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := m.EmailEditProfile(ctx, token, dbType, project, id, req["email"].(string), req["name"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}