package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandleSetCacheConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "cache-config", "modify", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		c := new(config.CacheConfig)
		_ = json.NewDecoder(r.Body).Decode(c)

		reqParams = utils.ExtractRequestParams(r, reqParams, c)
		status, err := syncMan.SetCacheConfig(ctx, c, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, status, w)
	}
}

// HandleGetCacheConfig is an endpoint handler which gets cache config
func HandleGetCacheConfig(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		reqParams, err := adminMan.IsTokenValid(ctx, token, "cache-config", "read", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		reqParams = utils.ExtractRequestParams(r, reqParams, nil)
		status, cacheConfig, err := syncMan.GetCacheConfig(ctx, reqParams)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, status, model.Response{Result: []interface{}{cacheConfig}})
	}
}

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandleGetCacheConnectionState(adminMan *admin.Manager, caching *caching.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "cache-config", "read", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: caching.ConnectionState(ctx)})
	}
}

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandlePurgeCache(adminMan *admin.Manager, caching *caching.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		c := new(model.CachePurgeRequest)
		_ = json.NewDecoder(r.Body).Decode(c)

		// Check if the request is authorised
		_, err := adminMan.IsTokenValid(ctx, token, "cache-config", "read", map[string]string{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		if err := caching.PurgeCache(ctx, projectID, c); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, model.Response{Result: caching.ConnectionState(ctx)})
	}
}

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandleInstantInvalidate(modules *modules.Modules) http.HandlerFunc {
	type request struct {
		Type string `json:"type"`
		Data struct {
			DB  string                 `json:"db"`
			Col string                 `json:"col"`
			Doc map[string]interface{} `json:"doc"`
		} `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := r.Header.Get("x-sc-token")
		arr := strings.Split(token, " ")
		if strings.ToLower(arr[0]) == "bearer" {
			token = arr[1]
		}

		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(utils.DefaultContextTime)*time.Second)
		defer cancel()

		// get project id from url
		vars := mux.Vars(r)
		projectID := vars["project"]

		c := new(request)
		_ = json.NewDecoder(r.Body).Decode(c)

		// Check if the request is authorised
		auth, err := modules.Auth(projectID)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		if err := auth.IsSCAccessToken(ctx, token); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		if err := modules.GlobalMods.Caching().InvalidateDatabaseCache(ctx, projectID, c.Data.DB, c.Data.Col, c.Type, c.Data.Doc); err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err)
			return
		}

		_ = helpers.Response.SendOkayResponse(ctx, http.StatusOK, w)
	}
}
