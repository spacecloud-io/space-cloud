package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleInspectTrackedCollectionsSchema is an endpoint handler which return schema for all tracked collections of a particular database
func HandleInspectTrackedCollectionsSchema(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Minute)
		defer cancel()

		vars := mux.Vars(r)
		dbAlias := vars["dbAlias"]
		projectID := vars["project"]

		schema := modules.Schema()
		schemas, err := schema.GetCollectionSchema(ctx, projectID, dbAlias)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: schemas})
	}
}
