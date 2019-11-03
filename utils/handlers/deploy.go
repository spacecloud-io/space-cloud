package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/modules/deploy"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
)

// HandleUploadAndDeploy handles the end to end deployment process
func HandleUploadAndDeploy(adminMan *admin.Manager, deploy *deploy.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, "deploy")
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := deploy.UploadAndDeploy(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
