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

// HandleUserManagement returns the handler to get the project config and validate the user via a REST endpoint
func HandleUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the body of the request
		value := new(config.AuthStub)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		provider := vars["provider"]

		// Sync the config
		if err := syncMan.SetUserManagement(ctx, projectID, provider, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give a positive acknowledgement
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

//GetUserManagement returns handler to get auth
func GetUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		providerID, exists := r.URL.Query()["provider"]

		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if exists {
			provider := project.Modules.Auth[providerID[0]]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"provider": provider})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"providers": project.Modules.Auth})
	}
}
