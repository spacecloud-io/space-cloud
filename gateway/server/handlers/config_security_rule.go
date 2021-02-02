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

// HandleSetSecurityFunction returns the handler to set security function
func HandleSetSecurityFunction(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := vars["id"]

		// Load the body of the request
		value := new(config.SecurityFunction)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := adminMan.IsTokenValid(ctx, token, "security-function", "modify", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		// Sync the config
		reqParams = utils.ExtractRequestParams(r, reqParams, value)
		status, err := syncMan.SetSecurityFunction(ctx, projectID, id, value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		// Give a positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetSecurityFunction returns handler to get security function
func HandleGetSecurityFunction(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		providerID := "*"
		providerQuery, exists := r.URL.Query()["id"]
		if exists {
			providerID = providerQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "security-function", "read", map[string]string{"project": projectID, "id": providerID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, providers, err := syncMan.GetSecurityFunction(ctx, projectID, providerID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: providers})
	}
}

// HandleDeleteSecurityFunction returns handler to delete security function
func HandleDeleteSecurityFunction(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		providerID := vars["id"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "security-function", "delete", map[string]string{"project": projectID, "id": providerID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, err := syncMan.DeleteSecurityFunction(ctx, projectID, providerID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}
		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}
