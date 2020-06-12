package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleSetFileSecretRootPath() http.HandlerFunc {
	type request struct {
		RootPath string `json:"rootPath"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)
		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		secretName := vars["id"]

		// Parse request body
		reqBody := new(request)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			logrus.Errorf("Failed to set file secret root path - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// set file secret root path
		if err := s.driver.SetFileSecretRootPath(ctx, projectID, secretName, reqBody.RootPath); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleApplySecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"]

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		secretObj.ID = name

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// create/update secret
		if err := s.driver.CreateSecret(ctx, projectID, secretObj); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleListSecrets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		name, exists := r.URL.Query()["id"]

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// list all secrets
		secrets, err := s.driver.ListSecrets(ctx, projectID)
		if err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
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

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"]

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// list all secrets
		if err := s.driver.DeleteSecret(ctx, projectID, name); err != nil {
			logrus.Errorf("Failed to delete secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleSetSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
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
			logrus.Errorf("Failed to set secret key - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// setSecretKey
		if err := s.driver.SetKey(ctx, projectID, name, key, secretVal); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleDeleteSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["id"] // secret-name
		key := vars["key"] // secret-key
		// setSecretKey

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		if err := s.driver.DeleteKey(ctx, projectID, name, key); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
