package handlers

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Get is a handler module to return a source registered in space-cloud
type Get struct {
	GVR schema.GroupVersionResource `json:"gvr"`
}

// CaddyModule returns the Caddy module information.
func (Get) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_get_handler",
		New: func() caddy.Module { return new(Get) },
	}
}

// ServeHTTP handles the http request
func (h *Get) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get the path params
	name := getName(r.URL.Path)
	if name == "" {
		vars := r.URL.Query()
		resp, err := configman.List(h.GVR, vars.Get("package"))
		if err != nil {
			return utils.SendErrorResponse(w, http.StatusInternalServerError, err)
		}

		return utils.SendResponse(w, http.StatusOK, resp)
	}

	resp, err := configman.Get(h.GVR, name)
	if err != nil {
		return utils.SendErrorResponse(w, http.StatusBadRequest, err)
	}

	return utils.SendResponse(w, http.StatusOK, resp)
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*Get)(nil)
