package handlers

import (
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandleSetCacheConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Caching module is disabled")
	}
}

// HandleGetCacheConfig is an endpoint handler which gets cache config
func HandleGetCacheConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, model.Response{Result: []interface{}{map[string]interface{}{"enabled": false}}})
	}
}

// HandleGetCacheConnectionState is an endpoint handler returns connection status
func HandleGetCacheConnectionState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, model.Response{Result: false})
	}
}
