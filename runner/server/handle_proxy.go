package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleWaitServices() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		vars := mux.Vars(r)
		project := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]
		// Wait for the service to scale up
		if err := s.driver.WaitForService(ctx, &model.Service{ProjectID: project, ID: serviceID, Version: version}); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err)
			return
		}
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleScaleUpService() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)
		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to scaleUp service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		project := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]
		// Instruct driver to scale up
		if err := s.driver.ScaleUp(ctx, project, serviceID, version); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
