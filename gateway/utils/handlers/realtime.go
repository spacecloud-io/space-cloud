package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

// HandleRealtimeEvent handles the request coming from the eventing module
func HandleRealtimeEvent(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Load the params from the body
		eventDoc := model.CloudEventPayload{}
		_ = json.NewDecoder(r.Body).Decode(&eventDoc)
		defer utils.CloseTheCloser(r.Body)

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := state.Auth.IsTokenInternal(token); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := state.Realtime.HandleRealtimeEvent(r.Context(), &eventDoc); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleRealtimeProcessRequest handles the request received from the realtime module. This is a request sent to every gateway
// instance in the cluster to propagate realtime changes
func HandleRealtimeProcessRequest(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the params from the body
		eventDoc := model.CloudEventPayload{}
		_ = json.NewDecoder(r.Body).Decode(&eventDoc)
		defer utils.CloseTheCloser(r.Body)

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := state.Auth.IsTokenInternal(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := state.Realtime.ProcessRealtimeRequests(&eventDoc); err != nil {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}
