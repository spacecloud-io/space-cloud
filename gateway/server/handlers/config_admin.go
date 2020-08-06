package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)
		clusterType, err := syncMan.GetClusterType(adminMan)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		isProd, plan, quotas, loginURL, clusterName, licenseRenewal, licenseKey, licenseValue := adminMan.LoadEnv()
		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{
			"isProd":       isProd,
			"plan":         plan,
			"quotas":       quotas,
			"version":      utils.BuildVersion,
			"licenseKey":   licenseKey,
			"licenseValue": licenseValue,
			"clusterName":  clusterName,
			"nextRenewal":  licenseRenewal,
			"clusterType":  clusterType,
			"loginURL":     loginURL,
		})
	}
}

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(adminMan *admin.Manager) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Key  string `json:"key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Check if the request is authorised
		status, token, err := adminMan.Login(ctx, req.User, req.Key)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": token})
	}
}

// HandleRefreshToken creates the refresh-token endpoint
func HandleRefreshToken(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminMan.GetCredentials()
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		newToken, err := adminMan.RefreshToken(token)
		if err != nil {
			logrus.Errorf("Error while refreshing token handleRefreshToken - %s ", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": newToken})
	}
}

// HandleGetPermissions returns the permission the authenticated user has
func HandleGetPermissions(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := adminMan.IsTokenValid(token, "config-permission", "read", nil)
		if err != nil {
			logrus.Errorf("Error while refreshing token handleRefreshToken - %s ", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header

		status, permissions, err := adminMan.GetPermissions(ctx, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		_ = utils.SendResponse(w, status, model.Response{Result: permissions})
	}
}

// HandleGenerateTokenForMissionControl handles the request of creating internal tokens
func HandleGenerateTokenForMissionControl(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := adminMan.IsTokenValid(token, "internal-token", "access", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("Error while refreshing token handleRefreshToken - %s ", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status, newToken, err := syncMan.GetTokenForMissionControl(ctx, projectID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		_ = utils.SendResponse(w, status, model.Response{Result: newToken})
	}
}
