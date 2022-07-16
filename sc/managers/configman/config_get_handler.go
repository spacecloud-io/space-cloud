package configman

import (
	"errors"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/middlewares"
	"github.com/spacecloud-io/space-cloud/utils"
)

// ConfigGetHandler is a module to create config GET handlers
type ConfigGetHandler struct {
	logger      *zap.Logger
	appLoader   loadApp
	store       *Store
	configTypes map[string]model.ConfigTypes
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

	store, err := ctx.App("config_store")
	if err != nil {
		return err
	}

	h.store = store.(*Store)

	// Load all the configuration types
	app, _ := ctx.App("configman")
	h.configTypes = app.(*ConfigMan).GetConfigTypes()
	return nil
}

// ServeHTTP handles the http request
func (h *ConfigGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get meta information
	op, module, typeName, resourceName, err := extractPathParams(r.URL.Path, r.Method)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Check if request is authenticated
	reqParams := middlewares.GetRequestParams(r)
	if !middlewares.IsRequestAuthenticated(reqParams, true) {
		h.logger.Error("Request has not been authenticated")
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, errors.New("user is not authenticated to make this request"))
	}

	// Get the type definition
	typeDef, err := loadConfigTypeDefinition(h.configTypes, module, typeName)
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
	if schemaErrors, err := typeDef.VerifyObject(resourceObj, op, false); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-get hooks if any
	if typeDef.Controller.PreGet != nil {
		if err := typeDef.Controller.PreGet(r.Context(), resourceObj.Meta, h.store); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// Put object in store
	if op == "single" {
		resources, err := h.store.GetResource(r.Context(), &resourceObj.Meta)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}

		// Invoke post-get hooks if any
		if typeDef.Controller.PostGet != nil {
			if err := typeDef.Controller.PostGet(r.Context(), &model.ListResourceObjects{List: []*model.ResourceObject{resources}}, h.store); err != nil {
				_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
				return nil
			}
		}

		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusOK, resources)
		return nil
	}
	resources, err := h.store.GetResources(r.Context(), &resourceObj.Meta)
	if err != nil {
		return err
	}

	// Invoke post-get hooks if any
	if typeDef.Controller.PostGet != nil {
		if err := typeDef.Controller.PostGet(r.Context(), resources, h.store); err != nil {
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
