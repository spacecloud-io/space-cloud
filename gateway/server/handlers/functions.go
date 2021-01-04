package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	authHelpers "github.com/spaceuptech/space-cloud/gateway/modules/auth/helpers"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleFunctionCall creates a functions request endpoint
func HandleFunctionCall(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		serviceID := vars["service"]
		function := vars["func"]

		auth := modules.Auth()
		functions := modules.Functions()

		// Load the params from the body
		req := model.FunctionsRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		timeOut, err := functions.GetEndpointContextTimeout(r.Context(), projectID, serviceID, function)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Create a new context
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeOut)*time.Second)
		defer cancel()

		actions, reqParams, err := auth.IsFuncCallAuthorised(ctx, projectID, serviceID, function, token, req.Params)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, req)

		status, result, err := functions.CallWithContext(ctx, serviceID, function, token, reqParams, req.Params)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Receieved error from service call (%s:%s)", serviceID, function), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = authHelpers.PostProcessMethod(ctx, auth.GetAESKey(), actions, result)

		_ = helpers.Response.SendResponse(ctx, w, status, result)
	}
}
