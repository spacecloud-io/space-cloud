package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/modules/deploy"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// HandleUploadAndDeploy handles the end to end deployment process
func HandleUploadAndDeploy(adminMan *admin.Manager, deploy *deploy.Module, projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := getToken(r)

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, "deploy")
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := deploy.UploadAndDeploy(r, projects); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
