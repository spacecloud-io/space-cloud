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

// HandleSetUserManagement returns the handler to get the project config and validate the user via a REST endpoint
func HandleSetUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		provider := vars["id"]

		// Load the body of the request
		value := new(config.AuthStub)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := adminMan.IsTokenValid(ctx, token, "auth-provider", "modify", map[string]string{"project": projectID, "id": provider})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Sync the config
		reqParams = utils.ExtractRequestParams(r, reqParams, value)
		status, err := syncMan.SetUserManagement(ctx, projectID, provider, value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		// Give a positive acknowledgement
		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetUserManagement returns handler to get auth
func HandleGetUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
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
		reqParams, err := adminMan.IsTokenValid(ctx, token, "auth-provider", "read", map[string]string{"project": projectID, "id": providerID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, providers, err := syncMan.GetUserManagement(ctx, projectID, providerID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: providers})
	}
}

// HandleDeleteUserManagement returns handler to delete auth
func HandleDeleteUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		providerID := vars["id"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "auth-provider", "delete", map[string]string{"project": projectID, "id": providerID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, err := syncMan.DeleteUserManagement(ctx, projectID, providerID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}
		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}
