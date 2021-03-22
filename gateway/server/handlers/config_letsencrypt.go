package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleLetsEncryptWhitelistedDomain handles the lets encrypt config request
func HandleLetsEncryptWhitelistedDomain(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := vars["id"]

		value := config.LetsEncrypt{}
		defer utils.CloseTheCloser(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		value.ID = id

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "letsencrypt", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, value)
		status, err := syncMan.SetProjectLetsEncryptDomains(ctx, projectID, &value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetEncryptWhitelistedDomain returns handler to get Encrypt White listed Domain
func HandleGetEncryptWhitelistedDomain(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "letsencrypt", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, le, err := syncMan.GetLetsEncryptConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: []interface{}{le}})
	}
}
