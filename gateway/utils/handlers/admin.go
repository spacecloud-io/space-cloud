package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleGetCredentials is an endpoint handler which gets username pass
func HandleGetCredentials(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: adminMan.GetCredentials()})
	}
}
