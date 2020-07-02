package handlers

import (
	"encoding/json"
	"net/http"

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
		if err := auth.IsTokenInternal(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		if err := realtime.HandleRealtimeEvent(r.Context(), &eventDoc); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleRealtimeProcessRequest handles the request received from the realtime module. This is a request sent to every gateway
// instance in the cluster to propagate realtime changes
func HandleRealtimeProcessRequest(modules *modules.Modules) http.HandlerFunc {
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
		if err := auth.IsTokenInternal(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		if err := realtime.ProcessRealtimeRequests(&eventDoc); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
