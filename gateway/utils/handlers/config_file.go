package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleSetFileStore set the file storage config
func HandleSetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := new(config.FileStore)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := syncMan.SetFileStore(ctx, projectID, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})

		// return
	}
}

// HandleGetFileState gets file state
func HandleGetFileState(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if projectConfig.Modules.FileStore.Enabled && projectConfig.Modules.FileStore.Conn != "" {
			w.WriteHeader(http.StatusOK) //http status code
			_ = json.NewEncoder(w).Encode(map[string]bool{"status": true})
		} else {
			w.WriteHeader(http.StatusOK) //http status code
			_ = json.NewEncoder(w).Encode(map[string]bool{"status": false})
		}

		// return
	}
}

// HandleSetFileRule sets file rule
func HandleSetFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := new(config.FileRule)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleName := vars["ruleName"]
		value.Name = ruleName

		if err := syncMan.SetFileRule(ctx, projectID, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})

		// return
	}
}

// HandleDeleteFileRule deletes file rule
func HandleDeleteFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := new(config.FileRule)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleName := vars["ruleName"]

		if err := syncMan.SetDeleteFileRule(ctx, projectID, ruleName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status code
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}
