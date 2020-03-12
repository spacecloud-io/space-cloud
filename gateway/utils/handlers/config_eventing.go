package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleAddEventingTriggerRule is an endpoint handler which adds a trigger rule to eventing
func HandleAddEventingTriggerRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
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
		ruleName := vars["triggerName"]
		projectID := vars["project"]

		if err := syncMan.SetEventingRule(ctx, projectID, ruleName, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

//HandleGetEventingTriggers returns handler to get event trigger
func HandleGetEventingTriggers(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get projectId and ruleName from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		ruleName, exists := r.URL.Query()["ruleName"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//check if ruleName exists
		if exists {
			rule, ok := project.Modules.Eventing.Rules[ruleName[0]]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("ruleName(%s) not present in state", ruleName[0])})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"rule": rule})
			return
		}

		if len(project.Modules.Eventing.Rules) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("eventing rules not present in state")})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"rules": project.Modules.Eventing.Rules})
	}
}

// HandleDeleteEventingTriggerRule is an endpoint handler which deletes a rule in eventing
func HandleDeleteEventingTriggerRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
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
		ruleName := vars["triggerName"]
		projectID := vars["project"]

		if err := syncMan.SetDeleteEventingRule(ctx, projectID, ruleName); err != nil {
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
		projectID := vars["project"]

		if err := syncMan.SetEventingConfig(ctx, projectID, c.DBType, c.Enabled); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		// return

	}
}

//HandleGetEventingConfig returns handler to get event config
func HandleGetEventingConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"dbAlias": project.Modules.Eventing.DBType, "enabled": project.Modules.Eventing.Enabled})
	}
}

// HandleSetEventingSchema is an endpoint handler which sets a schema in eventing
func HandleSetEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type schemaRequest struct {
		Schema string `json:"schema"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["type"]

		if err := syncMan.SetEventingSchema(ctx, projectID, evType, c.Schema); err != nil {
			logrus.Errorf("Failed to set eventing schema - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

//HandleGetEventingSchema returns handler to get event schema
func HandleGetEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id and type from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		typ, exists := r.URL.Query()["type"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//check id typ is present in url
		if exists {
			schema, ok := project.Modules.Eventing.Schemas[typ[0]]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("type(%s) not present in state", typ[0])})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"schema": schema})
			return
		}

		if len(project.Modules.Eventing.Schemas) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("eventing schema not present in state")})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"schemas": project.Modules.Eventing.Schemas})
	}
}

// HandleDeleteEventingSchema is an endpoint handler which deletes a schema in eventing
func HandleDeleteEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("Failed to validate token for delete eventing schema - %s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["type"]

		if err := syncMan.SetDeleteEventingSchema(ctx, projectID, evType); err != nil {
			logrus.Errorf("Failed to delete eventing schema - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleAddEventingSecurityRule is an endpoint handler which adds a security rule in eventing
func HandleAddEventingSecurityRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type setSecurityRules struct {
		Rule *config.Rule `json:"rule"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("Failed to validate token for set eventing rules - %s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := setSecurityRules{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["type"]

		if err := syncMan.SetEventingSecurityRules(ctx, projectID, evType, c.Rule); err != nil {
			logrus.Errorf("Failed to add eventing rules - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

//HandleGetEventingSecurityRules returns handler to get event security rules
func HandleGetEventingSecurityRules(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//get project id and type from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		typ, exists := r.URL.Query()["type"]

		//get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		//check if typ present in url
		if exists {
			rule, ok := project.Modules.Eventing.SecurityRules[typ[0]]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("type(%s) not present in state", typ[0])})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"securityRule": rule})
			return
		}
		if len(project.Modules.Eventing.SecurityRules) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint("eventing rules not present in state")})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"securityRules": project.Modules.Eventing.SecurityRules})
	}
}

// HandleDeleteEventingSecurityRule is an endpoint handler which deletes a security rule in eventing
func HandleDeleteEventingSecurityRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			logrus.Errorf("Failed to validate token for delete eventing rules - %s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["type"]

		if err := syncMan.SetDeleteEventingSecurityRules(ctx, projectID, evType); err != nil {
			logrus.Errorf("Failed to delete eventing rules - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK) //http status codee
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
