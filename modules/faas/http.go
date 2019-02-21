package faas

import (
	"encoding/json"
	"net/http"

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

		authObj, _ := auth.GetAuthObj(tokens[0])

		dataBytes, err := m.Request(engine, function, req.Timeout, map[string]interface{}{"auth": authObj, "params": req.Params})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		data := map[string]interface{}{}
		err = json.Unmarshal(dataBytes, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create the result to be sent back
		resultBytes, err := json.Marshal(map[string]interface{}{"result": data})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resultBytes)
	}
}
