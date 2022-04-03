package configman

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/managers/configman/connector"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// ConfigDeleteHandler is a module to create config Delete handlers
type ConfigDeleteHandler struct {
	logger    *zap.Logger
	appLoader loadApp
	store     connector.ConfigManConnector
}

// CaddyModule returns the Caddy module information.
func (ConfigDeleteHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_delete_handler",
		New: func() caddy.Module { return new(ConfigDeleteHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *ConfigDeleteHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	h.appLoader = ctx.App

	store, err := ctx.App("configman")
	if err != nil {
		return err
	}

	h.store = store.(*ConfigMan).Connectors
	return nil
}

// ServeHTTP handles the http request
func (h *ConfigDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get meta information
	op, module, typeName, resourceName, err := extractPathParams(r.URL.Path)
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

	resourceObj := new(model.ResourceObject)
	resourceObj.Meta.Module = module
	resourceObj.Meta.Type = typeName
	resourceObj.Meta.Name = resourceName
	resourceObj.Meta.Parents = utils.GetQueryParams(r.URL.Query())

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(resourceObj); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-delete hooks if any
	if err := applyHooks(r.Context(), module, typeDef, model.PhasePreDelete, h.appLoader, resourceObj); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Remove object from store
	if op == "single" {
		if err := h.store.DeleteResource(r.Context(), &resourceObj.Meta); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	} else {
		if err := h.store.DeleteResources(r.Context(), &resourceObj.Meta); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// Invoke pre-delete hooks if any
	if err := applyHooks(r.Context(), module, typeDef, model.PhasePostDelete, h.appLoader, resourceObj); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigDeleteHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigDeleteHandler)(nil)
