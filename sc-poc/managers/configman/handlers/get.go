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
	GVR  schema.GroupVersionResource `json:"gvr"`
	Name string                      `json:"name"`
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
	resp, err := configman.Get(h.GVR, h.Name)
	if err != nil {
		_ = utils.SendErrorResponse(w, 500, err)
	}

	utils.SendResponse(w, 200, resp)
	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*Get)(nil)
