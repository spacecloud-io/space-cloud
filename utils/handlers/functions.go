package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// HandleFunctionCall creates a Functions request endpoint
func HandleFunctionCall(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		service := vars["service"]
		function := vars["func"]

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Return if the functions module is not enabled
		if !state.Functions.IsEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Load the params from the body
		req := model.FunctionsRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		authObj, err := state.Auth.IsFuncCallAuthorised(project, service, function, token, req.Params)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		result, err := state.Functions.Call(service, function, authObj, req.Params, int(req.Timeout))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
