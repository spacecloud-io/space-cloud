package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/modules/userman"
	"github.com/spaceuptech/space-cloud/utils"
)

// HandleProfile returns the handler for fetching single user profile
func HandleProfile(userManagement *userman.Module) http.HandlerFunc {
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

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		status, result, err := userManagement.Profile(ctx, token, dbType, project, id)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"user": result})
	}
}

// HandleProfiles returns the handler for fetching all user profiles
func HandleProfiles(userManagement *userman.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		status, result, err := userManagement.Profiles(ctx, token, dbType, project)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignIn returns the handler for email sign in
func HandleEmailSignIn(userManagement *userman.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := userManagement.EmailSignIn(ctx, dbType, project, req["email"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailSignUp returns the handler for email sign up
func HandleEmailSignUp(userManagement *userman.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := userManagement.EmailSignUp(ctx, dbType, project, req["email"].(string), req["name"].(string), req["pass"].(string), req["role"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

// HandleEmailEditProfile returns the handler for edit profile
func HandleEmailEditProfile(userManagement *userman.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]
		id := vars["id"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		status, result, err := userManagement.EmailEditProfile(ctx, token, dbType, project, id, req["email"].(string), req["name"].(string), req["pass"].(string))

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}
