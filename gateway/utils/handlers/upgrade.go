package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleUpgrade(manager *syncman.Manager) http.HandlerFunc {
	type request struct {
		UserName   string `json:"userName"`
		Key        string `json:"key"`
		ClusterID  string `json:"clusterId"`
		ClusterKey string `json:"clusterKey"`
		Version    int    `json:"version"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)
		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		clusterType, err := manager.ConvertToEnterprise(ctx, req.UserName, req.Key, req.ClusterID, req.ClusterKey, req.Version)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: clusterType})
	}
}

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleVersionUpgrade(adminMan *admin.Manager, manager *syncman.Manager) http.HandlerFunc {
	type request struct {
		Version int `json:"version"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Println("version upgrade endpoint called")
		defer utils.CloseTheCloser(r.Body)
		if err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		if err := manager.VersionUpgrade(ctx, req.Version); err != nil {
			w.Header().Set("Content-Type", "application/json")
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
