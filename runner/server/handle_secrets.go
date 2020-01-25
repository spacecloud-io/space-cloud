package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleApplySecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// create/update secret
		if err := s.driver.CreateSecret(projectID, secretObj); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) handleListSecrets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		// list all secrets
		if _, err := s.driver.ListSecrets(projectID); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) handleDeleteSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["name"]

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to delete secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// list all secrets
		if err := s.driver.DeleteSecret(projectID, name); err != nil {
			logrus.Errorf("Failed to delete secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) handleSetSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// get nameSpace and secretKey from requestUrl!
		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["name"] //secret-name
		key := vars["key"]   //secret-key

		// body will only contain "value": secretValue (not-encoded!)
		secretVal := new(model.SecretValue)
		if err := json.NewDecoder(r.Body).Decode(secretVal); err != nil {
			logrus.Errorf("Failed to set secret key - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		// setSecretKey
		if err := s.driver.SetKey(projectID, secretVal, name, key); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) handleDeleteSecretKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		name := vars["name"] //secret-name
		key := vars["key"]   //secret-key
		// setSecretKey
		if err := s.driver.DeleteKey(projectID, name, key); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}
