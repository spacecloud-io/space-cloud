package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, value)
		status, err := syncMan.SetEventingRule(ctx, projectID, ruleName, &value, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-trigger", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, rules, err := syncMan.GetEventingTriggerRules(ctx, projectID, id, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: rules})
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-trigger", "modify", map[string]string{"project": projectID, "id": ruleName})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)
		status, err := syncMan.SetDeleteEventingRule(ctx, projectID, ruleName, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-config", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		c := new(config.Eventing)
		_ = json.NewDecoder(r.Body).Decode(c)

		utils.ExtractRequestParams(r, &reqParams, c)
		status, err := syncMan.SetEventingConfig(ctx, projectID, c.DBAlias, c.Enabled, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-config", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// get project config
		utils.ExtractRequestParams(r, &reqParams, nil)

		status, e, err := syncMan.GetEventingConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: []interface{}{e}})
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to validate token for set eventing schema", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		c := schemaRequest{}
		_ = json.NewDecoder(r.Body).Decode(&c)

		utils.ExtractRequestParams(r, &reqParams, c)
		status, err := syncMan.SetEventingSchema(ctx, projectID, evType, c.Schema, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-schema", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, schemas, err := syncMan.GetEventingSchema(ctx, projectID, id, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: schemas})
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-schema", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to validate token for delete eventing schema", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)
		status, err := syncMan.SetDeleteEventingSchema(ctx, projectID, evType, reqParams)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to delete eventing schema", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		c := new(config.Rule)
		_ = json.NewDecoder(r.Body).Decode(&c)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to validate token for set eventing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, c)
		status, err := syncMan.SetEventingSecurityRules(ctx, projectID, evType, c, reqParams)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to add eventing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-rule", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, securityRules, err := syncMan.GetEventingSecurityRules(ctx, projectID, id, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: securityRules})
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "eventing-rule", "modify", map[string]string{"project": projectID, "id": evType})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to validate token for delete eventing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)
		status, err := syncMan.SetDeleteEventingSecurityRules(ctx, projectID, evType, reqParams)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to delete eventing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}
