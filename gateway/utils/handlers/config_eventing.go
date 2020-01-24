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
		json.NewDecoder(r.Body).Decode(&value)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		if err := syncMan.SetEventingRule(ctx, project, ruleName, value); err != nil {
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
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		ruleName := vars["ruleName"]
		project := vars["project"]

		if err := syncMan.SetDeleteEventingRule(ctx, project, ruleName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}

// HandleSetEventingConfig is an endpoint handler which sets col and dytype in eventing according to body
func HandleSetEventingConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := config.Eventing{}
		json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetEventingConfig(ctx, project, c.DBType, c.Col, c.Enabled); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}

func HandleSetEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type schemaRequest struct {
		Type   string `json:"type"`
		Schema string `json:"schema"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaRequest{}
		json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetEventingSchema(ctx, project, c.Type, c.Schema); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}

func HandleDeleteEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type schemaDeleteRequest struct {
		Type string `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaDeleteRequest{}
		json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetDeleteEventingSchema(ctx, project, c.Type); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

func HandleSetEventingStrict(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type setStrict struct {
		Strict bool `json:"strict"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := setStrict{}
		json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		project := vars["project"]

		if err := syncMan.SetEventingStrict(ctx, project, c.Strict); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return

	}
}
