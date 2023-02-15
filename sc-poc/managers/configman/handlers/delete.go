package handlers

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Delete is a handler module to delete a registered source
type Delete struct {
	GVR  schema.GroupVersionResource `json:"gvr"`
	Name string                      `json:"name"`
}

// CaddyModule returns the Caddy module information.
func (Delete) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_delete_handler",
		New: func() caddy.Module { return new(Delete) },
	}
}

// ServeHTTP handles the http request
func (h *Delete) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	err := configman.Delete(h.GVR, h.Name)
	if err != nil {
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err)
	}

	utils.SendOkayResponse(w, http.StatusOK)
	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*Delete)(nil)
