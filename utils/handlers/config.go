package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// HandleLoadEnv returns the handler to load the projects via a REST endpoint
func HandleLoadEnv(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"isProd": adminMan.LoadEnv()})
	}
}

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		json.NewDecoder(r.Body).Decode(req)
		defer r.Body.Close()

		// Check if the request is authorised
		status, token, err := adminMan.Login(req.User, req.Pass)
		if err != nil {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		c := syncMan.GetGlobalConfig()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "projects": c.Projects})
	}
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

		// Load config from file
		c := syncMan.GetGlobalConfig()

		// Create a projects array
		projects := []*config.Project{}

		// Iterate over all projects
		for _, p := range c.Projects {
			// Add the project to the array if user has read access
			_, err := adminMan.IsAdminOpAuthorised(token, p.ID)
			if err == nil {
				projects = append(projects, p)
			}

			// Add an empty collections object is not present already
			for k, v := range p.Modules.Crud {
				if v.Collections == nil {
					p.Modules.Crud[k].Collections = map[string]*config.TableRule{}
				}
			}
		}

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

		// Sync the config
		if err := syncMan.SetProjectGlobalConfig(c); err != nil {
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

		// Sync the config
		if err := syncMan.SetProjectConfig(c); err != nil {
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
