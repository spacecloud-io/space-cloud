package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

// HandleGetCollectionSchemas is an endpoint handler which return schema for all the collection in the config.crud
func HandleGetCollectionSchemas(adminMan *admin.Manager, projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Minute)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		schemas, err := state.Schema.GetCollectionSchema(ctx, projectID, dbAlias)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) // http status codee
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(model.Response{Result: schemas})
	}
}
