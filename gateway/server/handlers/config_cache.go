package handlers

import (
	"net/http"

	"github.com/spaceuptech/helpers"
)

// HandleSetCacheConfig is an endpoint handler which sets cache config
func HandleSetCacheConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Not implemented in gateway")
	}
}

// HandleGetCacheConfig is an endpoint handler which gets cache config
func HandleGetCacheConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Not implemented in gateway")
	}
}

// HandleGetCacheConnectionState is an endpoint handler returns connection status
func HandleGetCacheConnectionState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Not implemented in gateway")
	}
}
