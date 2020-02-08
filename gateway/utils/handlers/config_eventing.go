package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleAddEventingRule is an endpoint handler which adds a rule to eventing
func HandleAddEventingRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := config.EventingRule{}
		_ = json.NewDecoder(r.Body).Decode(&value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		if err := syncMan.SetEventingRule(ctx, project, ruleName, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

// HandleDeleteEventingRule is an endpoint handler which deletes a rule in eventing
func HandleDeleteEventingRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		if err := syncMan.SetDeleteEventingRule(ctx, project, ruleName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

// HandleSetEventingConfig is an endpoint handler which sets col and dytype in eventing according to body
func HandleSetEventingConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := config.Eventing{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetEventingConfig(ctx, project, c.DBType, c.Col, c.Enabled); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

// HandleSetEventingSchema sets the schema for a custom event
func HandleSetEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type schemaRequest struct {
		Type   string `json:"type"`
		Schema string `json:"schema"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetEventingSchema(ctx, project, c.Type, c.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

// HandleDeleteEventingSchema deletes the schema for a custom event
func HandleDeleteEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type schemaDeleteRequest struct {
		Type string `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaDeleteRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetDeleteEventingSchema(ctx, project, c.Type); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return
	}
}
