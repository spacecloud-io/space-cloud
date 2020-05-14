package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleGetQuotas is an endpoint handler which number of projects & databases that can be created
func HandleGetQuotas(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // http status codee
		_ = json.NewEncoder(w).Encode(model.Response{Result: adminMan.GetQuotas()})
	}
}

// HandleGetCredentials is an endpoint handler which gets username pass
func HandleGetCredentials(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // http status code
		_ = json.NewEncoder(w).Encode(model.Response{Result: adminMan.GetCredentials()})
	}
}
