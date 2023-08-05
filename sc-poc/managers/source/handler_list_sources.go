package source

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/utils"
)

type ListSources struct{}

func (ListSources) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_list_sources_handler",
		New: func() caddy.Module { return new(ListSources) },
	}
}

func (h *ListSources) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	sourcesGVR := GetRegisteredSources()
	return utils.SendResponse(w, http.StatusOK, sourcesGVR)
}

var _ caddyhttp.MiddlewareHandler = (*ListSources)(nil)
