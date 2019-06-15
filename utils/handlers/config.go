package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/modules/auth"
)

// HandleAdminLogin creates the admin login endpoint
func HandleAdminLogin(auth *auth.Module) http.HandlerFunc {

	type Request struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Load the request from the body
		req := new(Request)
		json.NewDecoder(r.Body).Decode(req)
		defer r.Body.Close()

		status, token, err := auth.AdminLogin(req.User, req.Pass)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
	}
}
