package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleUpgrade(admin *admin.Manager, manager *syncman.Manager) http.HandlerFunc {
	type request struct {
		LicenseKey   string `json:"licenseKey"`
		LicenseValue string `json:"licenseValue"`
		ClusterName  string `json:"clusterName"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := admin.CheckIfAdmin(r.Context(), token); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		err := manager.ConvertToEnterprise(ctx, token, req.LicenseKey, req.LicenseValue, req.ClusterName)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Successfully upgraded gateway to enterprise`, nil)
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

func HandleDownGrade(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := admin.CheckIfAdmin(r.Context(), token); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if !admin.IsRegistered() {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, "Cannot remove license already running in open source mode")
			return
		}
		admin.ResetQuotas()

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Successfully removed license`, nil)
		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleRenewLicense(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := adminMan.CheckIfAdmin(r.Context(), token); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.RenewLicense(ctx, token); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}

// HandleSetOfflineLicense returns the handler to set offline licenses
func HandleSetOfflineLicense(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	type request struct {
		License string `json:"license"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		_ = json.NewDecoder(r.Body).Decode(req)
		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := adminMan.CheckIfAdmin(r.Context(), token); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetOfflineLicense(ctx, req.License); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
