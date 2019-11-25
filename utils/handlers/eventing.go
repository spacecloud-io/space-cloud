package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/eventing"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
)

// HandleProcessEvent processes a transmitted event
func HandleProcessEvent(adminMan *admin.Manager, eventing *eventing.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eventDocs := []*model.EventDocument{}
		json.NewDecoder(r.Body).Decode(&eventDocs)
		defer r.Body.Close()

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Process the incoming events
		eventing.ProcessTransmittedEvents(eventDocs)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}

}

// HandleQueueEvent creates a queue event endpoint
func HandleQueueEvent(adminMan *admin.Manager, eventing *eventing.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load the params from the body
		req := model.QueueEventRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Return if the eventing module is not enabled
		if !eventing.IsEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the path parameters
		// vars := mux.Vars(r)
		// project := vars["project"]

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := eventing.QueueEvent(ctx, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
