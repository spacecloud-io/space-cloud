package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleFunctionCall creates a functions request endpoint
func HandleFunctionCall(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		service := vars["service"]
		function := vars["func"]

		auth := modules.Auth()
		functions := modules.Functions()

		// Load the params from the body
		req := model.FunctionsRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)
		if req.Timeout == 0 {
			req.Timeout = 10 // set default context to 10 second
		}
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(req.Timeout)*time.Second)
		defer cancel()

		claims, err := auth.IsFuncCallAuthorised(ctx, projectID, service, function, token, req.Params)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		result, err := functions.CallWithContext(ctx, service, function, token, claims, req.Params)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"result": result})
	}
}
