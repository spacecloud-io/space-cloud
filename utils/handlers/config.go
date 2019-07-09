package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(adminMan *admin.Manager, syncMan *syncman.SyncManager) http.HandlerFunc {

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
		status, token, err := adminMan.Login(req.User, req.Pass)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		c := syncMan.GetGlobalConfig()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "projects": c.Projects})
	}
}

// HandleStoreConfig returns the handler to load the config via a REST endpoint
func HandleStoreConfig(adminMan *admin.Manager, syncMan *syncman.SyncManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := getToken(r)

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token)
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
			json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type - " + err.Error()})
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
func HandleLoadConfig(adminMan *admin.Manager, syncMan *syncman.SyncManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := getToken(r)

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load config from file
		c := syncMan.GetGlobalConfig()

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": c.Projects, "deploy": c.Deploy})
	}
}

// HandleDeleteConfig returns the handler to load the config via a REST endpoint
func HandleDeleteConfig(adminMan *admin.Manager, syncMan *syncman.SyncManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := getToken(r)

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err = syncMan.DeleteConfig(token, project)
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

// HandleStoreDeploymentConfig returns the handler to store the deployment config
func HandleStoreDeploymentConfig(adminMan *admin.Manager, syncMan *syncman.SyncManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := getToken(r)

		// Load the body of the request
		c := new(config.Deploy)
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer r.Body.Close()

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Set the deploy config
		if err := syncMan.SetDeployConfig(token, c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func getToken(r *http.Request) string {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}
	return strings.TrimPrefix(tokens[0], "Bearer ")
}
