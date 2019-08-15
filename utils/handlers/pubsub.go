package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/pubsub"
)

// HandlePublishCall publishes to pubsub
func HandlePublishCall(pubsub *pubsub.Module, auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		subject := vars["subject"]

		// Load the params from the body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		status, err := pubsub.Publish(project, token, subject, b)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
