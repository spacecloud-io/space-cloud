package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleLetsEncryptWhitelistedDomain handles the lets encrypt config request
func HandleLetsEncryptWhitelistedDomain(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := config.LetsEncrypt{}
		defer utils.CloseTheCloser(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := vars["id"]
		value.ID = id
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetProjectLetsEncryptDomains(ctx, projectID, value); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetEncryptWhitelistedDomain returns handler to get Encrypt White listed Domain
func HandleGetEncryptWhitelistedDomain(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		// get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: []interface{}{project.Modules.LetsEncrypt}})
	}
}
