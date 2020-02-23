package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleProcessEventResponse gets response for event
func HandleProcessEventResponse(adminMan *admin.Manager, eventing *eventing.Module) http.HandlerFunc {
	type request struct {
		BatchID  string      `json:"batchID"`
		Response interface{} `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			logrus.Errorf("error handling process event response eventing feature isn't enabled")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("error handling process event response token not valid - %s", err.Error())
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Process the incoming events
		eventing.SendEventResponse(req.BatchID, req.Response)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleProcessEvent processes a transmitted event
func HandleProcessEvent(adminMan *admin.Manager, eventing *eventing.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eventDocs := []*model.EventDocument{}
		_ = json.NewDecoder(r.Body).Decode(&eventDocs)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			logrus.Errorf("error handling process event request eventing feature isn't enabled")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("error handling process event request token not valid - %s", err.Error())
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Process the incoming events
		eventing.ProcessTransmittedEvents(eventDocs)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleQueueEvent creates a queue event endpoint
func HandleQueueEvent(eventing *eventing.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load the params from the body
		req := model.QueueEventRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			logrus.Errorf("error handling queue event request eventing feature isn't enabled")
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
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		if req.IsSynchronous {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
