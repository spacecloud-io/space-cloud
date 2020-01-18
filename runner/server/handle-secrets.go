package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleUpsertSecret() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)
		// get nameSpace from requestUrl!
		vars := mux.Vars(r)

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to create secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// create/update secret
		if err := s.driver.CreateSecret(secretObj, vars["projectID"]); err != nil {
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
		// get nameSpace from requestUrl!
		vars := mux.Vars(r)

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// list all secrets
		if _, err := s.driver.ListSecrets(secretObj, vars["projectID"]); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) handleDeleteSecrets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)
		// get nameSpace from requestUrl!
		vars := mux.Vars(r)

		// Parse request body
		secretObj := new(model.Secret)
		if err := json.NewDecoder(r.Body).Decode(secretObj); err != nil {
			logrus.Errorf("Failed to delete secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// list all secrets
		if err := s.driver.DelSecrets(secretObj, vars["projectID"]); err != nil {
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
		// get nameSpace and secretKey from requestUrl!
		vars := mux.Vars(r)

		// body will only contain "value": secretValue (not-encoded!)
		secretVal := new(model.SecretValue)
		if err := json.NewDecoder(r.Body).Decode(secretVal); err != nil {
			logrus.Errorf("Failed to set secret key - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		// setSecretKey
		if err := s.driver.SetKey(secretVal, vars["projectID"], vars["secretName"], vars["secretKey"]); err != nil {
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
		vars := mux.Vars(r)
		// setSecretKey
		if err := s.driver.DelKey(vars["projectID"], vars["secretName"], vars["secretKey"]); err != nil {
			logrus.Errorf("Failed to list secret - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}
