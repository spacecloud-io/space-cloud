package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

// HandleProfile returns the handler for fetching single user profile
func HandleProfile(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		defer utils.CloseTheCloser(r.Body)

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]
		id := vars["id"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		status, result, err := state.UserManagement.Profile(ctx, token, dbAlias, projectID, id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"user": result})
	}
}

// HandleProfiles returns the handler for fetching all user profiles
func HandleProfiles(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		status, result, err := state.UserManagement.Profiles(ctx, token, dbAlias, projectID)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignIn returns the handler for email sign in
func HandleEmailSignIn(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := map[string]interface{}{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		status, result, err := state.UserManagement.EmailSignIn(ctx, dbAlias, projectID, req["email"].(string), req["pass"].(string))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignUp returns the handler for email sign up
func HandleEmailSignUp(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := map[string]interface{}{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		status, result, err := state.UserManagement.EmailSignUp(ctx, dbAlias, projectID, req["email"].(string), req["name"].(string), req["pass"].(string), req["role"].(string))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailEditProfile returns the handler for edit profile
func HandleEmailEditProfile(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		dbAlias := vars["dbAlias"]
		id := vars["id"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := map[string]interface{}{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		status, result, err := state.UserManagement.EmailEditProfile(ctx, token, dbAlias, projectID, id, req["email"].(string), req["name"].(string), req["pass"].(string))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}
