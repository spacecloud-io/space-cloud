package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
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

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-config", "modify", map[string]string{"project": projectID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetFileStore(ctx, projectID, value); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetFileStore returns handler to get file store
func HandleGetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-config", "read", map[string]string{"project": projectID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get project config
		fileConfig, err := syncMan.GetFileStoreConfig(ctx, projectID)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: fileConfig})
	}
}

// HandleGetFileState gets file state
func HandleGetFileState(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-config", "read", map[string]string{"project": projectID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		file := modules.File()
		if err := file.GetState(ctx); err != nil {
			logrus.Errorf("error handling file get state got error - %s", err.Error())
			_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: false, Error: err.Error()})
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: true})
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

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-rule", "modify", map[string]string{"project": projectID, "id": ruleName}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetFileRule(ctx, projectID, ruleName, value); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
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

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-rule", "read", map[string]string{"project": projectID, "id": ruleID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// get project config
		fileRules, err := syncMan.GetFileStoreRules(ctx, projectID, ruleID)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: fileRules})
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

		value := new(config.FileRule)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "filestore-rule", "modify", map[string]string{"project": projectID, "id": ruleName}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetDeleteFileRule(ctx, projectID, ruleName); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
