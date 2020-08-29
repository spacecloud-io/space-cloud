package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create project", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		// Parse request body
		project := new(model.Project)
		if err := json.NewDecoder(r.Body).Decode(project); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create project", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		project.ID = projectID

		// Apply the service config
		if err := s.driver.CreateProject(ctx, project); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create project", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create project", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		// Apply the service config
		if err := s.driver.DeleteProject(ctx, projectID); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to create project", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]

		// Parse request body
		service := new(model.Service)
		if err := json.NewDecoder(r.Body).Decode(service); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		service.ProjectID = projectID
		service.ID = serviceID
		service.Version = version

		// TODO: Override the project id present in the service object with the one present in the token if user not admin

		// Apply the service config
		if err := s.driver.ApplyService(ctx, service); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func (s *Server) handleGetLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)
		// get query params
		vars := mux.Vars(r)
		projectID := vars["project"]

		taskID := r.URL.Query().Get("taskId")
		replicaID := r.URL.Query().Get("replicaId")
		if replicaID == "" {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, "replica id not provided in query param")
			return
		}
		_, isFollow := r.URL.Query()["follow"]

		helpers.Logger.LogDebug(helpers.GetRequestID(nil), "Get logs process started", map[string]interface{}{"projectId": projectID, "taskId": taskID, "replicaId": replicaID, "isFollow": isFollow})
		pipeReader, err := s.driver.GetLogs(r.Context(), isFollow, projectID, taskID, replicaID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(nil), "Failed to get service logs", err, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
			return
		}
		defer utils.CloseTheCloser(pipeReader)

		reader := bufio.NewReader(pipeReader)
		// implement http flusher
		flusher, ok := w.(http.Flusher)
		if !ok {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "expected http.ResponseWriter to be an http.Flusher")
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		for {
			select {
			case <-r.Context().Done():
				helpers.Logger.LogDebug(helpers.GetRequestID(nil), "Context deadline reached for client request", map[string]interface{}{})
				return
			default:
				str, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF && !isFollow {
						helpers.Logger.LogDebug(helpers.GetRequestID(nil), "End of file reached for logs", map[string]interface{}{})
						return
					}
					helpers.Logger.LogDebug(helpers.GetRequestID(nil), "error occured while reading from pipe in hander", nil)
					_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
					return
				}
				// starting 8 bytes of data contains some meta data regarding each log that docker sends
				// ignoring the first 8 bytes, send rest of the data
				fmt.Fprint(w, str)
				// Trigger "chunked" encoding and send a chunk...
				flusher.Flush()
			}
		}
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]
		version := vars["version"]

		if err := s.driver.DeleteService(ctx, projectID, serviceID, version); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID, serviceIDExists := r.URL.Query()["serviceId"]
		version, versionExists := r.URL.Query()["version"]

		services, err := s.driver.GetServices(ctx, projectID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
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

// HandleGetServicesStatus handles the request to get all services status
func (s *Server) HandleGetServicesStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		//var result []interface{}
		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID, serviceIDExists := r.URL.Query()["serviceId"]
		version, versionExists := r.URL.Query()["version"]

		result, err := s.driver.GetServiceStatus(ctx, projectID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to get service status", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		arr := make([]interface{}, 0)
		if serviceIDExists && versionExists {
			for _, serviceStatus := range result {
				if serviceStatus.ServiceID == serviceID[0] && serviceStatus.Version == version[0] {
					arr = append(arr, serviceStatus)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.Response{Result: arr})
			return
		}

		if serviceIDExists {
			for _, serviceStatus := range result {
				if serviceStatus.ServiceID == serviceID[0] {
					arr = append(arr, serviceStatus)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.Response{Result: arr})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Response{Result: result})
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to set service routes", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["serviceId"]

		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)

		err = s.driver.ApplyServiceRoutes(ctx, projectID, serviceID, req.Routes)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to apply service routing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to get service routes", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID, exists := r.URL.Query()["id"]

		serviceRoutes, err := s.driver.GetServiceRoutes(ctx, projectID)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to get service routing rules", err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
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
		ctx := r.Context()

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

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Proxy is making request to host (%s) port (%s)", ogHost, ogPort), nil)

		// Add to active request count
		// TODO: add support for multiple versions
		s.chAppend <- &model.ProxyMessage{Service: service, Project: project, Version: ogVersion, NodeID: "s-proxy", ActiveRequests: 1}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Wait for the service to scale up
		if err := s.debounce.Wait(fmt.Sprintf("proxy-%s-%s", project, service), func() error {
			return s.driver.WaitForService(ctx, &model.Service{ProjectID: project, ID: service, Version: ogVersion})
		}); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err.Error())
			return
		}

		var res *http.Response
		for i := 0; i < 5; i++ {
			// Fire the request
			var err error
			res, err = http.DefaultClient.Do(r)
			if err != nil {
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
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
