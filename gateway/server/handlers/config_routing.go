package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleSetProjectRoute adds a route in specified project config
func HandleSetProjectRoute(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := vars["id"]

		value := new(config.Route)
		_ = json.NewDecoder(r.Body).Decode(value)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()
		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "ingress-route", "modify", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling set project route in handlers unable to validate token got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, value)
		status, err := syncMan.SetProjectRoute(ctx, projectID, id, value, reqParams)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling set project route in handlers unable to add route in project config got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetProjectRoute returns handler to get project route
func HandleGetProjectRoute(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// get project id and routes id from url
		vars := mux.Vars(r)
		projectID := vars["project"]
		routeID := "*"
		routesQuery, exists := r.URL.Query()["id"]
		if exists {
			routeID = routesQuery[0]
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()
		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "ingress-route", "read", map[string]string{"project": projectID, "id": routeID})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, routes, err := syncMan.GetIngressRouting(ctx, projectID, routeID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: routes})
	}
}

// HandleDeleteProjectRoute deletes the specified route from project config
func HandleDeleteProjectRoute(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		vars := mux.Vars(r)
		projectID := vars["project"]
		routeID := vars["id"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		reqParams, err := adminMan.IsTokenValid(ctx, token, "ingress-route", "modify", map[string]string{"project": projectID, "id": routeID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling delete project route in handlers unable to validate token got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)
		status, err := syncMan.DeleteProjectRoute(ctx, projectID, routeID, reqParams)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling delete project route in handlers unable to delete route in project config got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleSetGlobalRouteConfig sets the project level ingress route config
func HandleSetGlobalRouteConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the required path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

		config := new(config.GlobalRoutesConfig)
		_ = json.NewDecoder(r.Body).Decode(config)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, utils.GetTokenFromHeader(r), "ingress-global", "modify", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling delete project route in handlers unable to validate token got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, config)
		status, err := syncMan.SetGlobalRouteConfig(ctx, projectID, config, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		// Send an okay response
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleGetGlobalRouteConfig gets the project level ingress route config
func HandleGetGlobalRouteConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the required path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, utils.GetTokenFromHeader(r), "ingress-global", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error handling delete project route in handlers unable to validate token got error message", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.ExtractRequestParams(r, &reqParams, nil)

		status, c, err := syncMan.GetGlobalRouteConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		// Send the repsonse back
		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: []interface{}{c}})
	}
}
