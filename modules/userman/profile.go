package userman

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// HandleProfile returns the handler for fetching single user profile
func (m *Module) HandleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow this feature only if the user management is enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]
		id := vars["id"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the user is authicated
		authObj, err := m.auth.IsAuthenticated(token, dbType, "users", utils.Read)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
			return
		}

		// Create the find object
		find := map[string]interface{}{}

		switch utils.DBType(dbType) {
		case utils.Mongo:
			find["_id"] = id

		default:
			find["id"] = id
		}

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"find": find, "op": utils.One, "auth": authObj},
			"project": project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = m.auth.IsAuthorized(dbType, "users", utils.Read, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		req := &model.ReadRequest{Find: find, Operation: utils.One}
		user, err := m.crud.Read(ctx, dbType, project, "users", req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
	}
}

// HandleProfiles returns the handler for fetching all user profiles
func (m *Module) HandleProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow this feature only if the user management is enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Check if the user is authicated
		authObj, err := m.auth.IsAuthenticated(token, dbType, "users", utils.Read)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authenticated"})
			return
		}

		// Create the find object
		find := map[string]interface{}{}

		// Create an args object
		args := map[string]interface{}{
			"args":    map[string]interface{}{"find": find, "op": utils.All, "auth": authObj},
			"project": project, // Don't forget to do this for every request
		}

		// Check if user is authorized to make this request
		err = m.auth.IsAuthorized(dbType, "users", utils.Read, args)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		req := &model.ReadRequest{Find: find, Operation: utils.All}
		users, err := m.crud.Read(ctx, dbType, project, "users", req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
	}
}
