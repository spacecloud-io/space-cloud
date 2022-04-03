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

// ConfigDeleteHandler is a module to create config Delete handlers
type ConfigDeleteHandler struct {
	logger    *zap.Logger
	appLoader loadApp

	store *ConfigMan
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

	h.store = store.(*ConfigMan)

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
	hook, err := loadHook(module, typeDef, model.PhasePreDelete, h.appLoader)
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
		if err := h.store.Connectors.DeleteResource(r.Context(), resourceMeta); err != nil {
			return err
		}
	} else {
		if err := h.store.Connectors.DeleteResources(r.Context(), resourceMeta); err != nil {
			return err
		}
	}

	// Invoke post-apply hooks if any
	hook, err = loadHook(module, typeDef, model.PhasePostDelete, h.appLoader)
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

	_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)

	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigDeleteHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigDeleteHandler)(nil)
