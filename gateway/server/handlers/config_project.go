package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

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

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(token, "project", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		status, project, err := syncMan.GetProjectConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		_ = utils.SendResponse(w, status, model.Response{Result: project})
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

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(token, "project", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = projectConfig
		statusCode, err := syncman.ApplyProjectConfig(ctx, &projectConfig, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, statusCode, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleDeleteProjectConfig returns the handler to delete the config of a project via a REST endpoint
func HandleDeleteProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Give negative acknowledgement
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, "Operation not supported")
	}
}
