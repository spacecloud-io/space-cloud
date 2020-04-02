package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/modules"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
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
			logrus.Errorf("error handling process event response eventing feature isn't enabled")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := auth.IsTokenInternal(token); err != nil {
			logrus.Errorf("error handling process event response token not valid - %s", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Process the incoming events
		eventing.SendEventResponse(req.BatchID, req.Response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
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
			logrus.Errorf("error handling process event request eventing feature isn't enabled")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("error handling process event request token not valid - %s", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Process the incoming events
		eventing.ProcessTransmittedEvents(eventDocs)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
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
			logrus.Errorf("error handling queue event request eventing feature isn't enabled")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
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
			logrus.Errorf("error handling queue event request - %s", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if res != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
