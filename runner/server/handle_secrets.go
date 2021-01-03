package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleSetFileSecretRootPath() http.HandlerFunc {
	type request struct {
		RootPath string `json:"rootPath"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)
		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		secretName := vars["id"]

		// Parse request body
		reqBody := new(request)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to set file secret root patt", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		// set file secret root path
		if err := s.driver.SetFileSecretRootPath(ctx, projectID, secretName, reqBody.RootPath); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleApplySecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"]

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		secretObj.ID = name

		// create/update secret
		if err := s.driver.CreateSecret(ctx, projectID, secretObj); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleListSecrets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		name, exists := r.URL.Query()["id"]

		// list all secrets
		secrets, err := s.driver.ListSecrets(ctx, projectID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to list secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		if exists {
			for _, val := range secrets {
				if val.ID == name[0] {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(model.Response{Result: []interface{}{val}})
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("secret(%s) not present in state", name[0])})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: secrets})
	}
}

func (s *Server) handleDeleteSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"]

		// list all secrets
		if err := s.driver.DeleteSecret(ctx, projectID, name); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to delete secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleSetSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace and secretKey from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"] // secret-name
		key := vars["key"] // secret-key

		// body will only contain "value": secretValue (not-encoded!)
		secretVal := new(model.SecretValue)
		if err := json.NewDecoder(r.Body).Decode(secretVal); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to set secret ket", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		// setSecretKey
		if err := s.driver.SetKey(ctx, projectID, name, key, secretVal); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to list secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleDeleteSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply servict", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"] // secret-name
		key := vars["key"] // secret-key
		// setSecretKey

		if err := s.driver.DeleteKey(ctx, projectID, name, key); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to list secret", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
