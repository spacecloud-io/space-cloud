package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (s *Server) handleCreateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		// Parse request body
		project := new(model.Project)
		if err := json.NewDecoder(r.Body).Decode(project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		project.ID = projectID

		// Apply the service config
		if err := s.driver.CreateProject(ctx, project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleDeleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		// Apply the service config
		if err := s.driver.DeleteProject(ctx, projectID); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

func (s *Server) handleApplyService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]

		// Parse request body
		service := new(model.Service)
		if err := json.NewDecoder(r.Body).Decode(service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		service.ProjectID = projectID
		service.ID = serviceID
		service.Version = version

		// TODO: Override the project id present in the service object with the one present in the token if user not admin

		// Apply the service config
		if err := s.driver.ApplyService(ctx, service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleDeleteService handles the request to delete a service
func (s *Server) HandleDeleteService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]

		if err := s.driver.DeleteService(ctx, projectID, serviceID, version); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetServices handles the request to get all services
func (s *Server) HandleGetServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID, serviceIDExists := r.URL.Query()["serviceId"]
		version, versionExists := r.URL.Query()["version"]

		services, err := s.driver.GetServices(ctx, projectID)
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result []*model.Service
		if serviceIDExists && versionExists {
			for _, val := range services {
				if val.ProjectID == projectID && val.ID == serviceID[0] && val.Version == version[0] {
					result = append(result, val)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(model.Response{Result: result})
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("serviceID(%s) or version(%s) not present in state", serviceID[0], version[0])})
			return
		}

		if serviceIDExists && !versionExists {
			for _, val := range services {
				if val.ID == serviceID[0] {
					result = append(result, val)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.Response{Result: result})
			return
		}

		result = services

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: result})
	}
}

// HandleApplyEventingService handles request to apply eventing service
func (s *Server) HandleApplyEventingService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		req := new(model.CloudEventPayload)
		_ = json.NewDecoder(r.Body).Decode(req)

		if req.Data.Meta.IsDeploy {
			// verify path e.g -> /artifacts/acc_id/projectid/version/build.zip
			arr := strings.Split(req.Data.Path, "/")
			// 7 will ensure that there will not be any index out of range error
			if len(arr) != 7 || arr[3] != req.Data.Meta.Service.ProjectID || arr[4] != req.Data.Meta.Service.Version {
				logrus.Errorf("error applying service path verification failed")
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, "error applying service path verification failed")
				return
			}
			// Apply the service config
			if err := s.driver.ApplyService(ctx, req.Data.Meta.Service); err != nil {
				logrus.Errorf("Failed to apply service - %s", err.Error())
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleServiceRoutingRequest handles request to apply service routing rules
func (s *Server) HandleServiceRoutingRequest() http.HandlerFunc {
	type request struct {
		Routes model.Routes `json:"routes"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to set service routes - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]

		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)

		err = s.driver.ApplyServiceRoutes(ctx, projectID, serviceID, req.Routes)
		if err != nil {
			logrus.Errorf("Failed to apply service routing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleGetServiceRoutingRequest handles request to get all service routing rules
func (s *Server) HandleGetServiceRoutingRequest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to get service routes - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID, exists := r.URL.Query()["id"]

		serviceRoutes, err := s.driver.GetServiceRoutes(ctx, projectID)
		if err != nil {
			logrus.Errorf("Failed to get service routing rules - %s", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if exists {
			for k, result := range serviceRoutes {
				if k == serviceID[0] {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(model.Response{Result: result})
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("serviceID(%s) not present in state", serviceID[0])})
			return
		}

		var result model.Routes
		for _, value := range serviceRoutes {
			result = append(result, value...)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: result})
	}
}

func (s *Server) handleProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.RequestURI = ""

		// Get the meta data from headers
		project := r.Header.Get("x-og-project")
		service := r.Header.Get("x-og-service")
		ogHost := r.Header.Get("x-og-host")
		ogPort := r.Header.Get("x-og-port")
		ogVersion := r.Header.Get("x-og-version")

		// Delete the headers
		r.Header.Del("x-og-project")
		r.Header.Del("x-og-service")
		r.Header.Del("x-og-host")
		r.Header.Del("x-og-port")
		r.Header.Del("x-og-version")

		// Change the destination with the original host and port
		r.Host = ogHost
		r.URL.Host = fmt.Sprintf("%s:%s", ogHost, ogPort)

		// Set the url scheme to http
		r.URL.Scheme = "http"

		logrus.Debugf("Proxy is making request to host (%s) port (%s)", ogHost, ogPort)

		// Add to active request count
		// TODO: add support for multiple versions
		s.chAppend <- &model.ProxyMessage{Service: service, Project: project, Version: ogVersion, NodeID: "s-proxy", ActiveRequests: 1}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Wait for the service to scale up
		if err := s.debounce.Wait(fmt.Sprintf("proxy-%s-%s", project, service), func() error {
			return s.driver.WaitForService(ctx, &model.Service{ProjectID: project, ID: service, Version: ogVersion})
		}); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		var res *http.Response
		for i := 0; i < 5; i++ {
			// Fire the request
			var err error
			res, err = http.DefaultClient.Do(r)
			if err != nil {
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			// TODO: Make this retry logic better
			if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusServiceUnavailable {
				break
			}

			time.Sleep(350 * time.Millisecond)

			// Close the body
			_, _ = io.Copy(ioutil.Discard, res.Body)
			utils.CloseTheCloser(res.Body)
		}

		defer utils.CloseTheCloser(res.Body)

		// Copy headers and status code
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}

		w.WriteHeader(res.StatusCode)
		_, _ = io.Copy(w, res.Body)
	}
}
