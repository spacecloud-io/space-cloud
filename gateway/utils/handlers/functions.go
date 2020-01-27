package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleFunctionCall creates a functions request endpoint
func HandleFunctionCall(functions *functions.Module, auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		service := vars["service"]
		function := vars["func"]

		// Load the params from the body
		req := model.FunctionsRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(req.Timeout)*time.Second)
		defer cancel()

		_, err := auth.IsFuncCallAuthorised(ctx, project, service, function, token, req.Params)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		result, err := functions.CallWithContext(ctx, service, function, token, req.Params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	}
}
