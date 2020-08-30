package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleSetFileStore set the file storage config
func HandleSetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		value := new(config.FileStore)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "filestore-config", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, value)
		status, err := syncMan.SetFileStore(ctx, projectID, value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetFileStore returns handler to get file store
func HandleGetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "filestore-config", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, fileConfig, err := syncMan.GetFileStoreConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: fileConfig})
	}
}

// HandleGetFileState gets file state
func HandleGetFileState(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		projectID := vars["project"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "filestore-config", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		file, err := modules.File(projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := file.GetState(ctx); err != nil {
			_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: false, Error: err.Error()})
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: true})
	}
}

// HandleSetFileRule sets file rule
func HandleSetFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleName := vars["id"]

		value := new(config.FileRule)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "filestore-rule", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, value)
		status, err := syncMan.SetFileRule(ctx, projectID, ruleName, value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetFileRule returns handler to get file rule
func HandleGetFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and ruleName
		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleID := "*"
		ruleName, exists := r.URL.Query()["id"]
		if exists {
			ruleID = ruleName[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "filestore-rule", "read", map[string]string{"project": projectID, "id": ruleID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		status, fileRules, err := syncMan.GetFileStoreRules(ctx, projectID, ruleID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: fileRules})
	}
}

// HandleDeleteFileRule deletes file rule
func HandleDeleteFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleName := vars["id"]

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "filestore-rule", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		if status, err := syncMan.SetDeleteFileRule(ctx, projectID, ruleName, reqParams); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
