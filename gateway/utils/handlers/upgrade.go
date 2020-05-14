package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleUpgrade(admin *admin.Manager, manager *syncman.Manager) http.HandlerFunc {
	type request struct {
		ClusterID  string `json:"clusterId"`
		ClusterKey string `json:"clusterKey"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := admin.IsTokenValid(token); err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		err := manager.ConvertToEnterprise(ctx, token, req.ClusterID, req.ClusterKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.LogDebug(`Successfully upgraded gateway to enterprise`, nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleRenewLicense(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := adminMan.IsTokenValid(token); err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.RenewLicense(ctx, token); err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{})
	}
}
