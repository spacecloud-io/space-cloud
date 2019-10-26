package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// HandleAddEventingRule is an endpoint handler which adds a rule to eventing
func HandleAddEventingRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		value := config.EventingRule{}
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		conf, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncMan.SetEventingRule(conf, ruleName, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}

// HandleDeleteEventingRule is an endpoint handler which deletes a rule in eventing
func HandleDeleteEventingRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncMan.SetDeleteEventingRule(projectConfig, ruleName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}

// HandleGetEventingStatus is an endpoint handler which sets col and dytype in eventing according to body
func HandleSetEventingStatus(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		type temp struct {
			Enabled bool   `json:"enabled" yaml:"enabled"`
			DBType  string `json:"db" yaml:"db"`
			Col     string `json:"col" yaml:"col"`
		}
		c := temp{}
		json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]
		conf, err := syncMan.GetConfig(project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := syncMan.SetEventingStatus(conf, c.DBType, c.Col, c.Enabled); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}
