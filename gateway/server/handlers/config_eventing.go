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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = value
		if err := syncMan.SetEventingRule(ctx, projectID, ruleName, value, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-trigger", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		rules, err := syncMan.GetEventingTriggerRules(ctx, projectID, id, reqParams)
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		if err := syncMan.SetDeleteEventingRule(ctx, projectID, ruleName, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-config", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := new(config.Eventing)
		_ = json.NewDecoder(r.Body).Decode(c)

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = c
		if err := syncMan.SetEventingConfig(ctx, projectID, c.DBAlias, c.Enabled, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-config", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get project config
		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		e, err := syncMan.GetEventingConfig(projectID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: []interface{}{e}})
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			logrus.Errorf("Failed to validate token for set eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := schemaRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = c
		if err := syncMan.SetEventingSchema(ctx, projectID, evType, c.Schema, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-schema", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		schemas, err := syncMan.GetEventingSchema(ctx, projectID, id, reqParams)
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			logrus.Errorf("Failed to validate token for delete eventing schema - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		if err := syncMan.SetDeleteEventingSchema(ctx, projectID, evType, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			logrus.Errorf("Failed to validate token for set eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		c := new(config.Rule)
		_ = json.NewDecoder(r.Body).Decode(&c)

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = c
		if err := syncMan.SetEventingSecurityRules(ctx, projectID, evType, c, reqParams); err != nil {
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-rule", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		securityRules, err := syncMan.GetEventingSecurityRules(ctx, projectID, id, reqParams)
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
		reqParams, err := adminMan.IsTokenValid(token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			logrus.Errorf("Failed to validate token for delete eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		if err := syncMan.SetDeleteEventingSecurityRules(ctx, projectID, evType, reqParams); err != nil {
			logrus.Errorf("Failed to delete eventing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
