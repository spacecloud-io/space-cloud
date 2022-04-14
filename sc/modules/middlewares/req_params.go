package middlewares

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/segmentio/ksuid"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
)

// RequestParamsHandler is responsible to add a request param object in the http.Request context
type RequestParamsHandler struct{}

// CaddyModule returns the Caddy module information.
func (RequestParamsHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_req_params_handler",
		New: func() caddy.Module { return new(RequestParamsHandler) },
	}
}

// ServeHTTP handles the http request
func (h *RequestParamsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Add a request id in the header if not already present
	if r.Header.Get(helpers.HeaderRequestID) == "" {
		r.Header.Set(helpers.HeaderRequestID, ksuid.New().String())
	}

	// Set the requset params in the context
	r = StoreRequestParams(r, utils.GenerateRequestParams(r))
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*RequestParamsHandler)(nil)
