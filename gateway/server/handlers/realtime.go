package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleRealtimeEvent handles the request coming from the eventing module
func HandleRealtimeEvent(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		realtime := modules.Realtime()

		// Load the params from the body
		eventDoc := model.CloudEventPayload{}
		_ = json.NewDecoder(r.Body).Decode(&eventDoc)
		defer utils.CloseTheCloser(r.Body)

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := auth.IsTokenInternal(r.Context(), token); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err)
			return
		}

		if err := realtime.HandleRealtimeEvent(r.Context(), &eventDoc); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	}
}
