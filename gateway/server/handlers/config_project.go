package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleGetProjectConfig returns handler to get config of the project
func HandleGetProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "project", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, project, err := syncMan.GetProjectConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: project})
	}
}

// HandleApplyProject is an endpoint handler which adds a project configuration in config
func HandleApplyProject(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		projectConfig := config.Project{}
		_ = json.NewDecoder(r.Body).Decode(&projectConfig)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig.ID = projectID

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "project", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, projectConfig)

		statusCode, err := syncman.ApplyProjectConfig(ctx, &projectConfig, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, statusCode, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, statusCode, w)
	}
}

// HandleDeleteProjectConfig returns the handler to delete the config of a project via a REST endpoint
func HandleDeleteProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Give negative acknowledgement
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, "Operation not supported")
	}
}

// HandleGetClusterConfig returns handler to get cluster-config
func HandleGetClusterConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "cluster", "read", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, clusterConfig, err := syncMan.GetClusterConfig(ctx, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: clusterConfig})
	}
}

// HandleSetClusterConfig set cluster-config
func HandleSetClusterConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the request from the body
		req := new(config.ClusterConfig)
		err := json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Throw error if request was of incorrect type
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Admin Config was of invalid type - "+err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "cluster", "modify", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, req)

		// Sync the Adminconfig
		status, err := syncMan.SetClusterConfig(ctx, req, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, map[string]interface{}{})
	}
}
