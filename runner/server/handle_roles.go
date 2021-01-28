package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

// HandleSetServiceRole handles request to apply service role
func (s *Server) HandleSetServiceRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to set service roles", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		req := new(model.Role)
		_ = json.NewDecoder(r.Body).Decode(req)

		vars := mux.Vars(r)
		req.Project = vars["project"]
		req.Service = vars["serviceId"]
		req.ID = vars["roleId"]

		err = s.driver.ApplyServiceRole(ctx, req)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleGetServiceRoleRequest handles request to get all service role
func (s *Server) HandleGetServiceRoleRequest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to get service role", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := r.URL.Query().Get("serviceId")
		roleID := r.URL.Query().Get("roleId")

		serviceRole, err := s.driver.GetServiceRole(ctx, projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		result := make([]*model.Role, 0)
		if serviceID != "" && roleID != "" {
			for _, role := range serviceRole {
				if role.ID == roleID && role.Service == serviceID {
					result = append(result, role)
				}
			}
			_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: result})
			return
		}

		if serviceID != "" {
			for _, role := range serviceRole {
				if role.Service == serviceID {
					result = append(result, role)
				}
			}
			_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: result})
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: serviceRole})
	}
}

// HandleDeleteServiceRole handles the request to delete a service role
func (s *Server) HandleDeleteServiceRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to delete service role", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]
		id := vars["roleId"]

		if err := s.driver.DeleteServiceRole(ctx, projectID, serviceID, id); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
