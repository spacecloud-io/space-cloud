package handlers

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List is a handler module to return all sources registered
// in space-cloud of a particular source type
type List struct {
	GVR schema.GroupVersionResource `json:"gvr"`
}

// CaddyModule returns the Caddy module information.
func (List) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_list_handler",
		New: func() caddy.Module { return new(List) },
	}
}

// ServeHTTP handles the http request
func (h *List) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	resp, err := configman.List(h.GVR)
	if err != nil {
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err)
	}

	utils.SendResponse(w, http.StatusOK, resp)
	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*List)(nil)
