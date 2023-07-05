package source

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

type ListPlugins struct {
	Plugins []v1alpha1.HTTPPlugin `json:"plugins"`
	logger  *zap.Logger
}

func (ListPlugins) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_list_plugins_handler",
		New: func() caddy.Module { return new(ListPlugins) },
	}
}

func (h *ListPlugins) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	sourceManT, err := ctx.App("source")
	if err != nil {
		h.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*App)
	h.Plugins = sourceMan.GetPlugins()
	return nil
}

func (h *ListPlugins) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return utils.SendResponse(w, http.StatusOK, h.Plugins)
}

var (
	_ caddy.Provisioner           = (*ListPlugins)(nil)
	_ caddyhttp.MiddlewareHandler = (*ListPlugins)(nil)
)
