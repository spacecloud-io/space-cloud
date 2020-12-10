package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleBatchApplyConfig applies all the config at once
func HandleBatchApplyConfig(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		req := new(model.BatchSpecApplyRequest)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		if err := adminMan.CheckIfAdmin(ctx, token); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		for _, specObject := range req.Specs {
			if err := utils.ApplySpec(ctx, token, "http://localhost:4122", specObject); err != nil {
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err.Error())
				return
			}
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
