package configman

import (
	"encoding/json"
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

// ConfigApplyHandler is a module to create config POST handlers
type ConfigApplyHandler struct {
	// Internal stuff
	logger      *zap.Logger
	appLoader   loadApp
	store       *Store
	configTypes map[string]model.ConfigTypes
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
func (h *ConfigApplyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get meta information
	op, module, typeName, _, err := extractPathParams(r.URL.Path, r.Method)
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

	// Extract the resourceObject object
	resourceObject := new(model.ResourceObject)
	if err := json.NewDecoder(r.Body).Decode(resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}
	resourceObject.Meta.Module = module
	resourceObject.Meta.Type = typeName
	resourceObject.Meta.Parents = utils.GetQueryParams(r.URL.Query())

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(resourceObject, op, true); err != nil {
		_ = helpers.Response.SendResponse(r.Context(), w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Invoke pre-apply hooks if any
	if typeDef.Controller.PreApply != nil {
		if err := typeDef.Controller.PreApply(r.Context(), resourceObject, h.store); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// Put object in store
	if err := h.store.ApplyResource(r.Context(), resourceObject); err != nil {
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return nil
	}

	// Invoke post-apply hooks if any
	if typeDef.Controller.PostApply != nil {
		if err := typeDef.Controller.PostApply(r.Context(), resourceObject, h.store); err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, err)
			return nil
		}
	}

	// Send ok response to client
	_ = helpers.Response.SendOkayResponse(r.Context(), http.StatusOK, w)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*ConfigApplyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*ConfigApplyHandler)(nil)
