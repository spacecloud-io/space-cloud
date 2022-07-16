package configman

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/modules/middlewares"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// OperationHandler is a module to create operation POST handlers
type OperationHandler struct {
	Operation string `json:"operation,omitempty"`

	// Internal stuff
	logger         *zap.Logger
	operationTypes map[string]model.OperationTypes
}

// CaddyModule returns the Caddy module information.
func (OperationHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_operation_handler",
		New: func() caddy.Module { return new(OperationHandler) },
	}
}

// Provision runs as a prehook to the handler operation
func (h *OperationHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	// Acquire the lock
	controllerLock.Lock()
	defer controllerLock.Unlock()

	// Load all the configuration types
	app, _ := ctx.App("configman")
	h.operationTypes = app.(*ConfigMan).GetOperationTypes()
	return nil
}

// ServeHTTP handles the http request
func (h *OperationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Prepare a context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Get meta information
	op, module, typeName, _, err := extractPathParams(r.URL.Path, r.Method)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
		return nil
	}

	// Get the type definition
	typeDef, err := loadOperationTypeDefinition(h.operationTypes, module, typeName)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
		return nil
	}

	reqParams := middlewares.GetRequestParams(r)
	if typeDef.IsProtected && !middlewares.IsRequestAuthenticated(reqParams, true) {
		h.logger.Error("Request has not been authenticated")
		_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, errors.New("user is not authenticated to make this request"))
	}

	// Check if incoming request is of the correct method
	if typeDef.Method != r.Method {
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, fmt.Errorf("invalid http method provided for operation %s/%s", module, typeName))
		return nil
	}

	// Extract the resourceObject object
	resourceObject := new(model.ResourceObject)
	resourceObject.Meta.Module = module
	resourceObject.Meta.Type = typeName
	resourceObject.Meta.Parents = utils.GetQueryParams(r.URL.Query())

	// Extract request body if method type was post or put
	if utils.StringExists([]string{http.MethodPost, http.MethodPut}, typeDef.Method) {
		payload, err := typeDef.Controller.DecodePayload(ctx, r.Body)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, err)
			return nil
		}

		resourceObject.Spec = payload
	}

	// Verify config object
	if schemaErrors, err := typeDef.VerifyObject(resourceObject, op, true); err != nil {
		_ = helpers.Response.SendResponse(ctx, w, http.StatusBadRequest, prepareErrorResponseBody(err, schemaErrors))
		return nil
	}

	// Process the operation
	status, res, err := typeDef.Controller.Handle(ctx, resourceObject, reqParams)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
		return nil
	}

	// Send the response
	_ = helpers.Response.SendResponse(ctx, w, status, res)
	return nil
}

// Interface guard
var _ caddy.Provisioner = (*OperationHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*OperationHandler)(nil)
