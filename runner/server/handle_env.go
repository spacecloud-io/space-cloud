package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleGetClusterType() http.HandlerFunc {
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
		logrus.Println("driver", s.driver.Type())
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: s.driver.Type()})
	}
}
