package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/modules"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleEventResponse gets response for event
func HandleEventResponse(modules *modules.Modules) http.HandlerFunc {
	type request struct {
		BatchID  string      `json:"batchID"`
		Response interface{} `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		auth := modules.Auth()
		eventing := modules.Eventing()

		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling process event response eventing feature isn't enabled", nil, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusNotFound, "This feature isn't enabled")
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := auth.IsTokenInternal(r.Context(), token); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling process event response token not valid", err, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err.Error())
			return
		}

		// Process the incoming events
		eventing.SendEventResponse(r.Context(), req.BatchID, req.Response)

		_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	}
}

// HandleProcessEvent processes a transmitted event
func HandleProcessEvent(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eventing := modules.Eventing()

		eventDocs := []*model.EventDocument{}
		_ = json.NewDecoder(r.Body).Decode(&eventDocs)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling process event request eventing feature isn't enabled", nil, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusNotFound, "This feature isn't enabled")
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if _, err := adminMan.IsTokenValid(r.Context(), token, "eventing-process", "process", nil); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling process event request token not valid", err, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err.Error())
			return
		}

		// Process the incoming events
		eventing.ProcessTransmittedEvents(eventDocs)

		_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	}
}

// HandleQueueEvent creates a queue event endpoint
func HandleQueueEvent(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eventing := modules.Eventing()

		// Load the params from the body
		req := model.QueueEventRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling queue event request eventing feature isn't enabled", nil, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusNotFound, "This feature isn't enabled")
			return
		}

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		res, err := eventing.QueueEvent(ctx, projectID, token, &req)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling queue event request", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		if res != nil {
			_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"result": res})
			return
		}
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
