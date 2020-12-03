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

// HandleAdminQueueEvent creates a queue event endpoint
func HandleAdminQueueEvent(adminMan *admin.Manager, modules *modules.Modules) http.HandlerFunc {
	type request struct {
		Events []*model.QueueEventRequest `json:"events"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		eventing := modules.Eventing()

		// Load the params from the body
		req := request{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling queue event request eventing feature isn't enabled", nil, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusNotFound, "This feature isn't enabled")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Get the JWT token from header
		if err := adminMan.CheckIfAdmin(ctx, utils.GetTokenFromHeader(r)); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusForbidden, err.Error())
			return
		}

		// Queue the event
		if err := eventing.QueueAdminEvent(ctx, req.Events); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), "error handling queue event request", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleQueueEvent creates a queue event endpoint
func HandleQueueEvent(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

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
