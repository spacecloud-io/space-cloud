package handlers

import (
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// HandleCluster returns handler cluster
func HandleCluster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		_ = helpers.Response.SendResponse(r.Context(), w, 501, model.Response{Error: "not implemented in open source"})
	}
}
