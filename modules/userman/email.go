package userman

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// HandleEmailSignIn returns the handler for email sign in
func (m *Module) HandleEmailSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow this feature only if the email sign in function is enabled
		if !m.isActive("email") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email sign in feature is not enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create read request
		readReq := &model.ReadRequest{Find: map[string]interface{}{"email": req["email"], "pass": req["pass"]}, Operation: utils.One}

		user, err := m.crud.Read(ctx, dbType, project, "users", readReq)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}

		// Delete password from user
		userObj := user.(map[string]interface{})
		delete(userObj, "pass")
		delete(req, "pass")

		// Create a token
		req["role"] = userObj["role"]
		req["name"] = userObj["name"]
		token, err := m.auth.CreateToken(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create a JWT token"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"user": user, "token": token})
	}
}

// HandleEmailSignUp returns the handler for email sign up
func (m *Module) HandleEmailSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow this feature only if the email sign in function is enabled
		if !m.isActive("email") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email sign in feature is not enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		// Load the request from the body
		req := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Create read request
		readReq := &model.ReadRequest{Find: map[string]interface{}{"email": req["email"], "pass": req["pass"]}, Operation: utils.One}
		_, err := m.crud.Read(ctx, dbType, project, "users", readReq)
		if err == nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "User with provided email already exists"})
			return
		}

		// Create a create request
		id := uuid.NewV1()
		if dbType == "mongo" {
			req["_id"] = id.String()
		} else {
			req["id"] = id.String()
		}
		createReq := &model.CreateRequest{Operation: utils.One, Document: req}
		err = m.crud.Create(ctx, dbType, project, "users", createReq)
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user account"})
			return
		}

		delete(req, "pass")
		token, err := m.auth.CreateToken(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create a JWT token"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"user": req, "token": token})
	}
}
