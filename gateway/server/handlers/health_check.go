package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
)

// HandleHealthCheck check health of gateway
func HandleHealthCheck(syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := syncMan.HealthCheck(); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
