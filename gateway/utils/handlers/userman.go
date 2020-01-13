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
		defer r.Body.Close()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]
		id := vars["id"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		status, result, err := state.UserManagement.Profile(ctx, token, dbType, project, id)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"user": result})
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
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		status, result, err := state.UserManagement.Profiles(ctx, token, dbType, project)
		if err != nil {
			return
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignIn returns the handler for email sign in
func HandleEmailSignIn(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		dbType := vars["dbType"]
		project := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := state.UserManagement.EmailSignIn(ctx, dbType, project, req["email"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
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
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := state.UserManagement.EmailSignUp(ctx, dbType, project, req["email"].(string), req["name"].(string), req["pass"].(string), req["role"].(string))
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
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
		id := vars["id"]
		dbType := vars["dbType"]
		project := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := state.UserManagement.EmailEditProfile(ctx, token, dbType, project, id, req["email"].(string), req["name"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}
