package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/managers/configman"
	"github.com/spacecloud-io/space-cloud/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Apply is a handler module to create/update source.
type Apply struct {
	GVR schema.GroupVersionResource `json:"gvr"`
}

// CaddyModule returns the Caddy module information.
func (Apply) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_apply_handler",
		New: func() caddy.Module { return new(Apply) },
	}
}

// ServeHTTP handles the http request
func (h *Apply) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	spec := &unstructured.Unstructured{}
	if err := json.NewDecoder(r.Body).Decode(spec); err != nil {
		_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
		return nil
	}

	if err := configman.Apply(h.GVR, spec); err != nil {
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err)
	}

	utils.SendOkayResponse(w, http.StatusOK)
	return nil
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*Apply)(nil)
