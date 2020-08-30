package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
		if err := admin.CheckIfAdmin(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		err := manager.ConvertToEnterprise(ctx, token, req.LicenseKey, req.LicenseValue, req.ClusterName)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.LogDebug(`Successfully upgraded gateway to enterprise`, "syncman", "startOperation", nil)
		_ = utils.SendOkayResponse(w, http.StatusOK)
	}
}

func HandleDownGrade(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := admin.CheckIfAdmin(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		_, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if !admin.IsRegistered() {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, "Cannot remove license already running in open source mode")
			return
		}
		admin.ResetQuotas()

		utils.LogDebug(`Successfully removed license`, "syncman", "startOperation", nil)
		_ = utils.SendOkayResponse(w, http.StatusOK)
	}
}

// HandleUpgrade returns the handler to load the projects via a REST endpoint
func HandleRenewLicense(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer utils.CloseTheCloser(r.Body)

		token := utils.GetTokenFromHeader(r)
		if err := adminMan.CheckIfAdmin(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.RenewLicense(ctx, token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w, http.StatusOK)
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
		if err := adminMan.CheckIfAdmin(token); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := syncMan.SetOfflineLicense(ctx, req.License); err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = utils.SendOkayResponse(w, http.StatusOK)
	}
}
