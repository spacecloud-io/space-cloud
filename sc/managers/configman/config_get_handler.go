package configman

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// ConfigGetHandler is a module to create config GET handlers
type ConfigGetHandler struct {
	logger    *zap.Logger
	appLoader loadApp

	store *ConfigMan
}

// CaddyModule returns the Caddy module information.
func (ConfigGetHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_config_get_handler",
		New: func() caddy.Module { return new(ConfigGetHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *ConfigGetHandler) Provision(ctx caddy.Context) error {
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
func (h *ConfigGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
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

	// Extract the resourceMeta object
	resourceMeta := new(model.ResourceMeta)

	resourceMeta.Module = module
	resourceMeta.Type = typeName
	resourceMeta.Name = resourceName
	resourceMeta.Parents = utils.GetQueryParams(r.URL.Query())

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(&model.ResourceObject{Meta: *resourceMeta}); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-apply hooks if any
	hook, err := loadHook(module, typeDef, model.PhasePreGet, h.appLoader)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Invoke hook if exists
	if hook != nil {
		if err := hook.Hook(r.Context(), &model.ResourceObject{Meta: *resourceMeta}); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// TODO: Put object in store
	if op == "single" {
		resources, err := h.store.Connectors.GetResource(r.Context(), resourceMeta)
		if err != nil {
			return err
		}

		// Invoke post-apply hooks if any
		hook, err = loadHook(module, typeDef, model.PhasePostGet, h.appLoader)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}

		// Invoke hook if exists
		if hook != nil {
			if err := hook.Hook(r.Context(), &model.ResourceObject{Meta: *resourceMeta}); err != nil {
				_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
				return nil
			}
		}

		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, resources)
		return nil
	}
	resources, err := h.store.Connectors.GetResources(r.Context(), resourceMeta)
	if err != nil {
		return err
	}

	// Invoke post-apply hooks if any
	hook, err = loadHook(module, typeDef, model.PhasePostGet, h.appLoader)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Invoke hook if exists
	if hook != nil {
		if err := hook.Hook(r.Context(), &model.ResourceObject{Meta: *resourceMeta}); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, resources)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigGetHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigGetHandler)(nil)
