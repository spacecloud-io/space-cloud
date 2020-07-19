package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

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

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(token, "ingress-route", "modify", map[string]string{"project": projectID, "id": id})
		if err != nil {
			logrus.Errorf("error handling set project route in handlers unable to validate token got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		reqParams.Payload = value
		if err := syncMan.SetProjectRoute(ctx, projectID, id, value, reqParams); err != nil {
			logrus.Errorf("error handling set project route in handlers unable to add route in project config got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
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

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(token, "ingress-route", "read", map[string]string{"project": projectID, "id": routeID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		routes, err := syncMan.GetIngressRouting(ctx, projectID, routeID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: routes})
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

		reqParams, err := adminMan.IsTokenValid(token, "ingress-route", "modify", map[string]string{"project": projectID, "id": routeID})
		if err != nil {
			logrus.Errorf("error handling delete project route in handlers unable to validate token got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		if err := syncMan.DeleteProjectRoute(ctx, projectID, routeID, reqParams); err != nil {
			logrus.Errorf("error handling delete project route in handlers unable to delete route in project config got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleSetGlobalRouteConfig sets the project level ingress route config
func HandleSetGlobalRouteConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type request struct {
		Config *config.GlobalRoutesConfig `json:"config"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the required path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

		// Get request body
		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r), "ingress-global", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling delete project route in handlers unable to validate token got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a new context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Set the config
		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		if err := syncMan.SetGlobalRouteConfig(ctx, projectID, req.Config, reqParams); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Send an okay response
		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetGlobalRouteConfig gets the project level ingress route config
func HandleGetGlobalRouteConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the required path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(utils.GetTokenFromHeader(r), "ingress-global", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling delete project route in handlers unable to validate token got error message - %v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a new context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Get the config from state
		reqParams.Method = r.Method
		reqParams.Path = r.URL.Path
		reqParams.Headers = r.Header
		c, err := syncMan.GetGlobalRouteConfig(ctx, projectID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Send the repsonse back
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: []interface{}{c}})
	}
}
