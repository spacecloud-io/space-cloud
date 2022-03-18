package configman

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// ConfigApplyHandler is a module to create config POST handlers
type ConfigApplyHandler struct {
	logger    *zap.Logger
	appLoader loadApp
}

// CaddyModule returns the Caddy module information.
func (ConfigApplyHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_apply_handler",
		New: func() caddy.Module { return new(ConfigApplyHandler) },
	}
}

func (h *ConfigApplyHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	h.appLoader = ctx.App

	return nil
}

func (h *ConfigApplyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get meta information
	_, module, typeName, _, err := extractPathParams(r.URL.Path)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Get the type definition
	typeDef, err := loadTypeDefinition(module, typeName)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Extract the configObject object
	configObject := new(ResourceObject)
	if err := json.NewDecoder(r.Body).Decode(configObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}
	configObject.Meta.Module = module
	configObject.Meta.Type = typeName

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(configObject); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-apply hooks if any
	hook, err := loadHook(module, typeDef, PhasePreApply, h.appLoader)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Invoke hook if exists
	if hook != nil {
		if err := hook.Hook(r.Context(), configObject); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// TODO: Put object in store

	// Send ok response to client
	_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigApplyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigApplyHandler)(nil)
