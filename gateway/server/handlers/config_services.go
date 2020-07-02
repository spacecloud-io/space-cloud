package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleAddService is an endpoint handler which deletes a table in specified database
func HandleAddService(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		vars := mux.Vars(r)
		service := vars["id"]
		projectID := vars["project"]

		v := config.Service{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		defer utils.CloseTheCloser(r.Body)

		reqParams, err := adminMan.IsTokenValid(token, "remote-service", "modify", map[string]string{"project": projectID, "service": service})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Headers = r.Header
		if err := syncMan.SetService(ctx, projectID, service, &v, reqParams); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetService returns handler to get services of the project
func HandleGetService(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := "*"
		serviceQuery, ok := r.URL.Query()["id"]
		if ok {
			serviceID = serviceQuery[0]
		}

		reqParams, err := adminMan.IsTokenValid(token, "remote-service", "read", map[string]string{"project": projectID, "service": serviceID})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		reqParams.Headers = r.Header
		services, err := syncMan.GetServices(ctx, projectID, serviceID, reqParams)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: services})
	}
}

// HandleDeleteService is an endpoint handler which deletes a table in specified database
func HandleDeleteService(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		vars := mux.Vars(r)
		service := vars["id"]
		projectID := vars["project"]

		defer utils.CloseTheCloser(r.Body)

		reqParams, err := adminMan.IsTokenValid(token, "remote-service", "modify", map[string]string{"project": projectID, "service": service})
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		reqParams.Headers = r.Header
		if err := syncMan.DeleteService(ctx, projectID, service, reqParams); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}
