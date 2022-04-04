package configman

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// ConfigApplyHandler is a module to create config POST handlers
type ConfigApplyHandler struct {
	logger    *zap.Logger
	appLoader loadApp
	store     *ConfigMan
}

// CaddyModule returns the Caddy module information.
func (ConfigApplyHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_apply_handler",
		New: func() caddy.Module { return new(ConfigApplyHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *ConfigApplyHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	h.appLoader = ctx.App

	store, err := ctx.App("configman")
	if err != nil {
		return err
	}

	h.store = store.(*ConfigMan)
	return nil
}

// ServeHTTP handles the http request
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

	// Extract the resourceObject object
	resourceObject := new(model.ResourceObject)
	if err := json.NewDecoder(r.Body).Decode(resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}
	resourceObject.Meta.Module = module
	resourceObject.Meta.Type = typeName

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(resourceObject); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-apply hooks if any
	if err := applyHooks(r.Context(), module, typeDef, model.PhasePreApply, h.appLoader, resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Put object in store
	if err := h.store.connector.ApplyResource(r.Context(), resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Invoke post-apply hooks if any
	if err := applyHooks(r.Context(), module, typeDef, model.PhasePostApply, h.appLoader, resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Send ok response to client
	_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigApplyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigApplyHandler)(nil)
