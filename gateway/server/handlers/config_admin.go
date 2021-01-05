package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleGetCredentials is an endpoint handler which gets username pass
func HandleGetCredentials(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if _, err := adminMan.IsTokenValid(ctx, utils.GetTokenFromHeader(r), "creds", "read", nil); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to validate token for set eventing schem", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: adminMan.GetCredentials()})
	}
}

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)

		clusterType, err := syncMan.GetClusterType(ctx, adminMan)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		isProd, plan, quotas, loginURL, clusterName, licenseRenewal, licenseKey, licenseValue, sessionID, licenseMode, err := adminMan.LoadEnv()
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{
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
			"sessionId":    sessionID,
			"licenseMode":  licenseMode,
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
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)

		// Load the request from the body
		req := new(Request)
		_ = json.NewDecoder(r.Body).Decode(req)

		// Check if the request is authorised
		status, token, err := adminMan.Login(ctx, req.User, req.Key)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"token": token})
	}
}

// HandleRefreshToken creates the refresh-token endpoint
func HandleRefreshToken(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		newToken, err := adminMan.RefreshToken(ctx, token)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error while refreshing token handleRefreshToken", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"token": newToken})
	}
}

// HandleGetPermissions returns the permission the authenticated user has
func HandleGetPermissions(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		reqParams, err := adminMan.IsTokenValid(ctx, token, "config-permission", "read", nil)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error while refreshing token handleRefreshToken", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, permissions, err := adminMan.GetPermissions(ctx, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: permissions})
	}
}

// HandleGenerateTokenForMissionControl handles the request of creating internal tokens
func HandleGenerateTokenForMissionControl(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		defer utils.CloseTheCloser(r.Body)
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := adminMan.IsTokenValid(ctx, token, "internal-token", "access", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Error while refreshing token handleRefreshToken", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		status, newToken, err := syncMan.GetTokenForMissionControl(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: newToken})
	}
}

// HandleGenerateAdminToken generates an admin token with the claims provided
func HandleGenerateAdminToken(adminMan *admin.Manager) http.HandlerFunc {
	type Request struct {
		Claims map[string]interface{} `json:"claims"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := new(Request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		newToken, err := adminMan.GenerateToken(r.Context(), token, req.Claims)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err)
			return
		}

		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, model.Response{Result: map[string]string{"token": newToken}})
	}
}
