package server

import (
	"context"
	"encoding/json"
	"errors"
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
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// Parse request body
		project := new(model.Project)
		if err := json.NewDecoder(r.Body).Decode(project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// Apply the service config
		if err := s.driver.CreateProject(ctx, project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.SendEmptySuccessResponse(w, r)
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
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["projectId"]
		// Apply the service config
		if err := s.driver.DeleteProject(ctx, projectID); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.SendEmptySuccessResponse(w, r)
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
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		// Parse request body
		service := new(model.Service)
		if err := json.NewDecoder(r.Body).Decode(service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// TODO: Override the project id present in the service object with the one present in the token if user not admin

		// Apply the service config
		if err := s.driver.ApplyService(ctx, service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) HandleDeleteService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["projectId"]
		serviceID := vars["serviceId"]
		version := vars["version"]

		if err := s.driver.DeleteService(ctx, projectID, serviceID, version); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (s *Server) HandleGetServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["projectId"]

		services, err := s.driver.GetServices(ctx, projectID)
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"services": services})
	}
}

func (s *Server) HandleApplyEventingService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Verify token
		_, err := s.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		req := new(model.CloudEventPayload)
		json.NewDecoder(r.Body).Decode(req)

		if req.Data.Meta.IsDeploy {
			// verify path e.g -> /artifacts/acc_id/projectid/version/build.zip
			arr := strings.Split(req.Data.Path, "/")
			// 7 will ensure that there will not be any index out of range error
			if len(arr) != 7 || arr[3] != req.Data.Meta.Service.ProjectID || arr[4] != req.Data.Meta.Service.Version {
				logrus.Errorf("error applying service path verification failed")
				utils.SendErrorResponse(w, r, http.StatusInternalServerError, errors.New("error applying service path verification failed"))
				return
			}
			// Apply the service config
			if err := s.driver.ApplyService(ctx, req.Data.Meta.Service); err != nil {
				logrus.Errorf("Failed to apply service - %s", err.Error())
				utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		utils.SendEmptySuccessResponse(w, r)
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
		r.Header.Del("x-og-env")
		r.Header.Del("x-og-version")

		// Change the destination with the original host and port
		r.Host = ogHost
		r.URL.Host = fmt.Sprintf("%s:%s", ogHost, ogPort)

		// Set the url scheme to http
		r.URL.Scheme = "http"

		// Add to active request count
		// TODO: add support for multiple versions
		s.chAppend <- &model.ProxyMessage{Service: service, Project: project, Version: ogVersion, NodeID: "s-proxy", ActiveRequests: 1}

		// Wait for the service to scale up
		if err := s.debounce.Wait(fmt.Sprintf("proxy-%s-%s", project, service), func() error {
			return s.driver.WaitForService(&model.Service{ProjectID: project, ID: service, Version: ogVersion})
		}); err != nil {
			utils.SendErrorResponse(w, r, http.StatusServiceUnavailable, err)
			return
		}

		var res *http.Response
		for i := 0; i < 5; i++ {
			// Fire the request
			var err error
			res, err = http.DefaultClient.Do(r)
			if err != nil {
				utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
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
		w.WriteHeader(res.StatusCode)
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}

		_, _ = io.Copy(w, res.Body)
	}
}
