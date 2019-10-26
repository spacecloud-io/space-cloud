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

// HandleSetFileStore set the file storage config
func HandleSetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		value := new(config.FileStore)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncMan.SetFileStore(projectConfig, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})

		return
	}
}

// HandleGetFileState gets file state
func HandleGetFileState(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if projectConfig.Modules.FileStore.Enabled && projectConfig.Modules.FileStore.Conn != "" {
			w.WriteHeader(http.StatusOK) //http status code
			json.NewEncoder(w).Encode(map[string]bool{"status": true})
		} else {
			w.WriteHeader(http.StatusOK) //http status code
			json.NewEncoder(w).Encode(map[string]bool{"status": false})
		}

		return
	}
}

// HandleSetFileRule sets file rule
func HandleSetFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		value := new(config.FileRule)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]
		ruleName := vars["ruleName"]
		value.Name = ruleName

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		projectConfig.Modules.FileStore.Rules = append(projectConfig.Modules.FileStore.Rules, value)

		if err := syncMan.SetFileRule(projectConfig, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})

		return
	}
}

// HandleDeleteFileRule deletes file rule
func HandleDeleteFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		value := new(config.FileRule)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]
		filename := vars["filename"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncMan.SetDeleteFileRule(projectConfig, filename); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}
