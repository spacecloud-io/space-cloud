package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"isProd": adminMan.LoadEnv(), "version": utils.BuildVersion})
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
		json.NewDecoder(r.Body).Decode(req)
		defer r.Body.Close()

		// Check if the request is authorised
		status, token, err := adminMan.Login(req.User, req.Key)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		c := syncMan.GetGlobalConfig()

		token, err = adminMan.GetInternalAccessToken()
		if err != nil {
			logrus.Error("error in admin login of handler unable to generate internal access token - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		for _, project := range c.Projects {
			services, err := getServices(syncMan, project.ID, token)
			if err != nil {
				logrus.Error("error in admin login of handler unable to set deployments - %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			project.Modules.Deployments.Services = services
		}
		syncMan.SetGlobalConfig(c)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "projects": c.Projects})
	}
}

func getServices(syncMan *syncman.Manager, projectID, token string) ([]*config.RunnerService, error) {
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

	type resp struct {
		Services []*config.RunnerService `json:"services"`
		Error    string                  `json:"error"`
	}
	data := resp{}
	if err = json.NewDecoder(httpRes.Body).Decode(&data); err != nil {
		logrus.Errorf("error while getting services in handler unable to decode response boyd -%v", err)
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		logrus.Errorf("error while getting services in handler got http request -%v", httpRes.StatusCode)
		return nil, fmt.Errorf("error while getting services in handler got http request -%v -%v", httpRes.StatusCode, data.Error)
	}

	return data.Services, err
}

// HandleLoadProjects returns the handler to load the projects via a REST endpoint
func HandleLoadProjects(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		adminToken, err := adminMan.GetInternalAccessToken()
		if err != nil {
			logrus.Error("error while loading projects handlers unable to generate internal access token - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load config from file
		c := syncMan.GetGlobalConfig()

		// Create a projects array
		projects := []*config.Project{}

		// Iterate over all projects
		for _, p := range c.Projects {
			// Add the project to the array if user has read access
			_, err := adminMan.IsAdminOpAuthorised(token, p.ID)
			if err == nil {
				services, err := getServices(syncMan, p.ID, adminToken)
				if err != nil {
					logrus.Error("error while loading projects in handler unable to get services - %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
					return
				}
				p.Modules.Deployments.Services = services
				projects = append(projects, p)
			}

			// Add an empty collections object is not present already
			for k, v := range p.Modules.Crud {
				if v.Collections == nil {
					p.Modules.Crud[k].Collections = map[string]*config.TableRule{}
				}
			}
		}
		syncMan.SetGlobalConfig(c)

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": projects})
	}
}

// HandleGlobalConfig returns the handler to store the global config of a project via a REST endpoint
func HandleGlobalConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Load the body of the request
		c := new(config.Project)
		err := json.NewDecoder(r.Body).Decode(c)
		defer r.Body.Close()

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type - " + err.Error()})
			return
		}

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, c.ID)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		if err := syncMan.SetProjectGlobalConfig(ctx, c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
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
		defer r.Body.Close()

		// Throw error if request was of incorrect type
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Config was of invalid type - " + err.Error()})
			return
		}

		// Check if the request is authorised
		status, err := adminMan.IsAdminOpAuthorised(token, c.ID)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Sync the config
		if err := syncMan.SetProjectConfig(ctx, c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// HandleDeleteProjectConfig returns the handler to delete the config of a project via a REST endpoint
func HandleDeleteProjectConfig(adminMan *admin.Manager, syncMan *syncman.Manager, configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		// Give negative acknowledgement
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Operation not supported"})
	}
}
