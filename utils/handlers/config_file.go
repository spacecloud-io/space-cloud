package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

func SetFileStore(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		value := new(config.FileStore)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {

		}
		projectConfig.Modules.FileStore.Enabled = value.Enabled
		projectConfig.Modules.FileStore.StoreType = value.StoreType
		projectConfig.Modules.FileStore.Conn = value.Conn
		projectConfig.Modules.FileStore.Endpoint = value.Endpoint
		projectConfig.Modules.FileStore.Bucket = value.Bucket

		syncMan.SetProjectConfig(projectConfig)

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})

		return
	}
}

func GetFileState(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		vars := mux.Vars(r)
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {

		}

		if projectConfig.Modules.FileStore.Enabled && projectConfig.Modules.FileStore.Conn != "" {
			w.WriteHeader(http.StatusOK) //http status code
			json.NewEncoder(w).Encode(map[string]interface{}{"status": true})
		} else {
			w.WriteHeader(http.StatusOK) //http status code
			json.NewEncoder(w).Encode(map[string]interface{}{"status": false})
		}

		return
	}
}

func SetFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		value := new(config.FileRule)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {

		}
		projectConfig.Modules.FileStore.Rules = append(projectConfig.Modules.FileStore.Rules, value)

		syncMan.SetProjectConfig(projectConfig)

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})

		return
	}
}

func DeleteFileRule(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		value := new(config.FileRule)
		json.NewDecoder(r.Body).Decode(&value)

		vars := mux.Vars(r)
		project := vars["project"]
		filename := vars["filename"]

		projectConfig, err := syncMan.GetConfig(project)
		if err != nil {

		}
		temp := projectConfig.Modules.FileStore.Rules
		for i, v := range projectConfig.Modules.FileStore.Rules {
			if v.Name == filename {
				temp = append(temp[:i], temp[i+1:]...)
			}
		}
		projectConfig.Modules.FileStore.Rules = temp

		syncMan.SetProjectConfig(projectConfig)

		w.WriteHeader(http.StatusOK) //http status code
		json.NewEncoder(w).Encode(map[string]interface{}{})

		return
	}
}
