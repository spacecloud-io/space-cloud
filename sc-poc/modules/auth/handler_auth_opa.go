package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/provider"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

// OPAPlugin is responsible to run OPA policy on the incomming request
type OPAPlugin struct {
	Name string `json:"name"`
	Rego string `json:"rego"`

	logger      *zap.Logger
	providerMan *provider.App
}

// CaddyModule returns the Caddy module information.
func (OPAPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_plugin_opa_handler",
		New: func() caddy.Module { return new(OPAPlugin) },
	}
}

// Provision sets up the opa handler module.
func (h *OPAPlugin) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	providerManT, err := ctx.App("provider")
	if err != nil {
		return err
	}
	h.providerMan = providerManT.(*provider.App)
	return nil
}

// ServeHTTP handles the http request
func (h *OPAPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get the auth app
	ws := source.GetWorkspaceNameFromHeaders(r)
	app, _ := h.providerMan.GetProvider(ws, "auth")
	authApp := app.(*App)

	// Get auth result
	authResult, p := utils.GetAuthenticationResult(r.Context())
	if !p {
		h.logger.Error("Unable to load authentication result")
		_ = utils.SendErrorResponse(w, http.StatusUnauthorized, errors.New("unable to authenticate request"))
		return nil
	}

	// Extract the body
	// TODO: Optimize this by extracting all request parameters in some handler and storing it in the request context.
	// All handlers can then ready the request parameters from the context instead of parsing it again and again from the body.
	body := map[string]interface{}{}
	if !utils.StringExists(r.Method, http.MethodGet, http.MethodDelete) {
		buffer := bytes.NewBuffer([]byte{})
		_ = json.NewDecoder(io.TeeReader(r.Body, buffer)).Decode(&body)

		// Reassign the body
		r.Body = io.NopCloser(buffer)
	}

	input := map[string]interface{}{
		"auth": map[string]interface{}{
			"isAuthenticated": authResult.IsAuthenticated,
			"claims":          authResult.Claims,
		},
		"request": map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"body":   body,
			// TODO: put query params
		},
	}

	// Calculate the policy name
	var policyName string
	if h.Name != "" {
		policyName = h.Name
	} else {
		policyName = utils.Hash(h.Rego)
	}

	allowed, reason, err := authApp.EvaluatePolicy(r.Context(), policyName, input)
	if err != nil {
		h.logger.Error("Unable to run opa policy", zap.Error(err))
		_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err)
		return nil
	}

	if !allowed {
		h.logger.Info("Opa policy failed",
			zap.String("policy", policyName), zap.String("reason", reason),
			zap.String("path", r.URL.Path), zap.String("method", r.Method))
		_ = utils.SendErrorResponse(w, http.StatusForbidden, errors.New(reason))
		return nil
	}

	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*OPAPlugin)(nil)
var _ caddyhttp.MiddlewareHandler = (*OPAPlugin)(nil)
