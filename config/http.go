package config

import (
	"encoding/json"
	"net/http"
)

// HandleConfig returns the handler to load the config via a REST endpoint
func HandleConfig(isProd bool, cb func(*Project) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Disable this feature if env is prod
		if isProd {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "This active only in the dev environment"})
			return
		}

		// Load the body of the request
		config := new(Project)
		err := json.NewDecoder(r.Body).Decode(config)
		defer r.Body.Close()

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type"})
			return
		}

		// Call the callback
		err = cb(config)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
