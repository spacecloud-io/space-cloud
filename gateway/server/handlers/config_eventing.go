package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleAddEventingTriggerRule is an endpoint handler which adds a trigger rule to eventing
func HandleAddEventingTriggerRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		ruleName := vars["id"]
		projectID := vars["project"]

		value := config.EventingRule{}
		_ = json.NewDecoder(r.Body).Decode(&value)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetEventingRule(ctx, projectID, ruleName, value); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetEventingTriggers returns handler to get event trigger
func HandleGetEventingTriggers(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get projectId and ruleName from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		id := "*"
		ruleNameQuery, exists := r.URL.Query()["id"]
		if exists {
			id = ruleNameQuery[0]
		}

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-trigger", "read", map[string]string{"project": projectID, "id": id}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		rules, err := syncMan.GetEventingTriggerRules(ctx, projectID, id)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: rules})
	}
}

// HandleDeleteEventingTriggerRule is an endpoint handler which deletes a rule in eventing
func HandleDeleteEventingTriggerRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		ruleName := vars["id"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetDeleteEventingRule(ctx, projectID, ruleName); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleSetEventingConfig is an endpoint handler which sets col and dytype in eventing according to body
func HandleSetEventingConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-config", "modify", map[string]string{"project": projectID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := new(config.Eventing)
		_ = json.NewDecoder(r.Body).Decode(c)
		if err := syncMan.SetEventingConfig(ctx, projectID, c.DBAlias, c.Enabled); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetEventingConfig returns handler to get event config
func HandleGetEventingConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-config", "read", map[string]string{"project": projectID}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get project config
		project, err := syncMan.GetConfig(projectID)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: []interface{}{config.Eventing{DBAlias: project.Modules.Eventing.DBAlias, Enabled: project.Modules.Eventing.Enabled}}})
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

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["id"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType}); err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		if err := syncMan.SetEventingSchema(ctx, projectID, evType, c.Schema); err != nil {
			logrus.Errorf("Failed to set eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetEventingSchema returns handler to get event schema
func HandleGetEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and type from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		id := "*"
		typ, exists := r.URL.Query()["id"]
		if exists {
			id = typ[0]
		}

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-schema", "read", map[string]string{"project": projectID, "id": id}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		schemas, err := syncMan.GetEventingSchema(ctx, projectID, id)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: schemas})
	}
}

// HandleDeleteEventingSchema is an endpoint handler which deletes a schema in eventing
func HandleDeleteEventingSchema(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["id"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType}); err != nil {
			logrus.Errorf("Failed to validate token for delete eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetDeleteEventingSchema(ctx, projectID, evType); err != nil {
			logrus.Errorf("Failed to delete eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleAddEventingSecurityRule is an endpoint handler which adds a security rule in eventing
func HandleAddEventingSecurityRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["id"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType}); err != nil {
			logrus.Errorf("Failed to validate token for set eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := new(config.Rule)
		_ = json.NewDecoder(r.Body).Decode(&c)
		if err := syncMan.SetEventingSecurityRules(ctx, projectID, evType, c); err != nil {
			logrus.Errorf("Failed to add eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetEventingSecurityRules returns handler to get event security rules
func HandleGetEventingSecurityRules(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and type from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		id := "*"
		typ, exists := r.URL.Query()["id"]
		if exists {
			id = typ[0]
		}

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-rule", "read", map[string]string{"project": projectID, "id": id}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		securityRules, err := syncMan.GetEventingSecurityRules(ctx, projectID, id)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: securityRules})
	}
}

// HandleDeleteEventingSecurityRule is an endpoint handler which deletes a security rule in eventing
func HandleDeleteEventingSecurityRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		projectID := vars["project"]
		evType := vars["id"]

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType}); err != nil {
			logrus.Errorf("Failed to validate token for delete eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetDeleteEventingSecurityRules(ctx, projectID, evType); err != nil {
			logrus.Errorf("Failed to delete eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
