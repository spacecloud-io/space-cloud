package handlers

import (
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleCluster returns handler cluster
func HandleCluster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		_ = utils.SendResponse(w, 501, model.Response{Error: "not implemented in open source"})
	}
}
