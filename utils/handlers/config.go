package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(auth *auth.Module) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		json.NewDecoder(r.Body).Decode(req)
		defer r.Body.Close()

		// Check if the request is authorised
		status, token, project, err := auth.AdminLogin(req.User, req.Pass)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "project": project})
	}
}

// HandleStoreConfig returns the handler to load the config via a REST endpoint
func HandleStoreConfig(auth *auth.Module, syncMan *syncman.SyncManager, configPath string, cb func(*config.Project) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		status, err := auth.IsAdminOpAuthorised(project, token)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the body of the request
		c := new(config.Project)
		err = json.NewDecoder(r.Body).Decode(c)
		defer r.Body.Close()

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type"})
			return
		}

		// Sync the config
		err = syncMan.SetConfig(token, c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleLoadConfig returns the handler to load the config via a REST endpoint
func HandleLoadConfig(auth *auth.Module, syncMan *syncman.SyncManager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		status, err := auth.IsAdminOpAuthorised(project, token)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load config from file
		c, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"config": c})
	}
}
