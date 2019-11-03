package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// HandleRealtimeEvent handles the request coming from the eventing module
func HandleRealtimeEvent(projects *projects.Projects) http.HandlerFunc {
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
		eventDoc := model.EventDocument{}
		json.NewDecoder(r.Body).Decode(&eventDoc)
		defer r.Body.Close()

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := state.Auth.IsTokenInternal(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := state.Realtime.HandleRealtimeEvent(r.Context(), &eventDoc); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

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
		eventDoc := model.EventDocument{}
		json.NewDecoder(r.Body).Decode(&eventDoc)
		defer r.Body.Close()

		// Get the token
		token := utils.GetTokenFromHeader(r)

		// Check if the token is valid
		if err := state.Auth.IsTokenInternal(token); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := state.Realtime.ProcessRealtimeRequests(&eventDoc); err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
