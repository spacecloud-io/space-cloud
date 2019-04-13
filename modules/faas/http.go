package faas

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"

	"github.com/spaceuptech/space-cloud/modules/auth"
)

// HandleRequest creates a FaaS request endpoint
func (m *Module) HandleRequest(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Return if the faas module is not enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Get the path parametrs
		vars := mux.Vars(r)
		engine := vars["engine"]
		function := vars["func"]

		// Load the params from the body
		req := model.FaaSRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		resultBytes, err := m.Operation(auth, token, engine, function, req.Timeout)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resultBytes)
	}
}
