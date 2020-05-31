package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)
		clusterType, err := syncMan.GetClusterType(adminMan)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"isProd": adminMan.LoadEnv(), "version": utils.BuildVersion, "plan": "space-cloud-open", "quotas": adminMan.GetQuotas(), "clusterId": "", "clusterType": clusterType})
	}
}

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Key  string `json:"key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		status, token, err := adminMan.Login(req.User, req.Key)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}

		c := syncMan.GetGlobalConfig()
		// if endpoint is called by cli don't insert deployments config in projects
		cli, ok := r.URL.Query()["cli"]
		if ok && cli[0] == "true" {
			_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": token, "projects": c.Projects})
			return
		}

		if syncMan.GetRunnerAddr() != "" {
			adminToken, err := adminMan.GetInternalAccessToken()
			if err != nil {
				logrus.Errorf("error while loading projects handlers unable to generate internal access token - %s", err.Error())
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			for _, project := range c.Projects {
				services, err := getServices(syncMan, project.ID, adminToken)
				if err != nil {
					logrus.Errorf("error in admin login of handler unable to set deployments - %s", err.Error())
					_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				project.Modules.Deployments.Services = services
				secrets, err := getSecrets(syncMan, project.ID, adminToken)
				if err != nil {
					logrus.Errorf("error in admin login of handler unable to set secrets - %s", err.Error())
					_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				project.Modules.Secrets = secrets
			}
			syncMan.SetGlobalConfig(c)
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": token, "projects": c.Projects})
	}
}

func getServices(syncMan *syncman.Manager, projectID, token string) (interface{}, error) {
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/v1/runner/%s/services", syncMan.GetRunnerAddr(), projectID), nil)
	if err != nil {
		logrus.Errorf("error while getting services in handler unable to create http request - %v", err)
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		logrus.Errorf("error while getting services in handler unable to execute http request - %v", err)
		return nil, err
	}

	data := new(model.Response)
	if err = json.NewDecoder(httpRes.Body).Decode(&data); err != nil {
		logrus.Errorf("error while getting services in handler unable to decode response body -%v", err)
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		logrus.Errorf("error while getting services in handler got http request -%v", httpRes.StatusCode)
		return nil, fmt.Errorf("error while getting services in handler got http request -%v -%v", httpRes.StatusCode, data.Error)
	}

	return data.Result, err
}

func getSecrets(syncMan *syncman.Manager, projectID, token string) (interface{}, error) {
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/v1/runner/%s/secrets", syncMan.GetRunnerAddr(), projectID), nil)
	if err != nil {
		logrus.Errorf("error while getting secrets in handler unable to create http request - %v", err)
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		logrus.Errorf("error while getting secrets in handler unable to execute http request - %v", err)
		return nil, err
	}

	data := new(model.Response)
	if err = json.NewDecoder(httpRes.Body).Decode(&data); err != nil {
		logrus.Errorf("error while getting secrets in handler unable to decode response body -%v", err)
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		logrus.Errorf("error while getting secrets in handler got http status code -%v", httpRes.StatusCode)
		return nil, fmt.Errorf("http status %v message -%v", httpRes.StatusCode, data.Error)
	}

	return data.Result, err
}

// HandleRefreshToken creates the refresh-token endpoint
func HandleRefreshToken(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		defer utils.CloseTheCloser(r.Body)

		newToken, err := adminMan.RefreshToken(token)
		if err != nil {
			logrus.Errorf("Error while refreshing token handleRefreshToken - %s ", err.Error())
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"token": newToken})
	}
}

// HandleLoadProjects returns the handler to load the projects via a REST endpoint
func HandleLoadProjects(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Load config from file
		c := syncMan.GetGlobalConfig()

		if syncMan.GetRunnerAddr() != "" {
			adminToken, err := adminMan.GetInternalAccessToken()
			if err != nil {
				logrus.Errorf("error while loading projects handlers unable to generate internal access token - %s", err.Error())
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			for _, project := range c.Projects {
				services, err := getServices(syncMan, project.ID, adminToken)
				if err != nil {
					logrus.Errorf("error in admin login of handler unable to set deployments - %s", err.Error())
					_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				project.Modules.Deployments.Services = services
				secrets, err := getSecrets(syncMan, project.ID, adminToken)
				if err != nil {
					logrus.Errorf("error in admin login of handler unable to set secrets - %s", err.Error())
					_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				project.Modules.Secrets = secrets
			}
			syncMan.SetGlobalConfig(c)
		}

		// Give positive acknowledgement
		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: c.Projects})
	}
}

// HandleGetProjectConfig returns handler to get config of the project
func HandleGetProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		project, err := syncMan.GetProjectConfig(projectID)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, model.Response{Result: project})
	}
}

// HandleApplyProject is an endpoint handler which adds a project configuration in config
func HandleApplyProject(adminMan *admin.Manager, syncman *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		projectConfig := config.Project{}
		_ = json.NewDecoder(r.Body).Decode(&projectConfig)
		defer utils.CloseTheCloser(r.Body)

		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig.ID = projectID

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		statusCode, err := syncman.ApplyProjectConfig(ctx, &projectConfig)
		if err != nil {
			_ = utils.SendErrorResponse(w, statusCode, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w)
	}
}

// HandleStoreProjectConfig returns the handler to store the config of a project via a REST endpoint
func HandleStoreProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the body of the request
		c := new(config.Project)
		err := json.NewDecoder(r.Body).Decode(c)
		defer utils.CloseTheCloser(r.Body)

		// Throw error if request was of incorrect type
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, "Config was of invalid type - "+err.Error())
			return
		}

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		if err := syncMan.SetProjectConfig(ctx, c); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Give positive acknowledgement
		_ = utils.SendOkayResponse(w)
	}
}

// HandleDeleteProjectConfig returns the handler to delete the config of a project via a REST endpoint
func HandleDeleteProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		// Give negative acknowledgement
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, "Operation not supported")
	}
}
