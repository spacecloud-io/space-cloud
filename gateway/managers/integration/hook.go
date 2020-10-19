package integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Manager) invokeHooks(ctx context.Context, params model.RequestParams) authResponse {
	var hookResponse authResponse

	// Extract the caller id if the caller is an integration
	var callerID string
	if role, p := params.Claims["role"]; p && role == "integration" {
		callerID = params.Claims["id"].(string)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Inside invoke hooks", map[string]interface{}{"resource": params.Resource, "verb": params.Op})
	// TODO: Make the iteration happen in parallel
	for _, hook := range m.integrationHookConfig {
		// Skip if the integration is the caller
		if callerID == hook.IntegrationID {
			continue
		}
		// Check if the resource types match
		if !utils.StringExists(hook.Resource, "*", params.Resource) {
			continue
		}

		// Check if the op matches
		if !utils.StringExists(hook.Verbs, "*", params.Op) {
			continue
		}

		// Check if the attr match
		if params.Attributes != nil {
			for k, v := range hook.Attributes {
				val, p := params.Attributes[k]
				if !p {
					continue
				}

				if !utils.StringExists(v, "*", val) {
					continue
				}
			}
		}

		// TODO: Check if rule matches

		// Get the admin token
		scToken, err := m.adminMan.GetInternalAccessToken()
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to make sc token for invoking hook (%s) in integration (%s)", hook.ID, hook.IntegrationID), err, nil)
			continue
		}

		// Invoke the hook
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Invoking hook", map[string]interface{}{
			"url":      hook.URL,
			"resource": params.Resource,
			"verb":     params.Op,
		})
		var res config.IntegrationHookResponse
		req := &utils.HTTPRequest{URL: hook.URL, Params: params, Method: http.MethodPost, SCToken: scToken}
		status, err := utils.MakeHTTPRequest(ctx, req, &res)
		if err != nil {
			err = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to make request to hook (%s) in integration (%s)", hook.ID, hook.IntegrationID), err, nil)

			// Return error if this was a hook with hijack permission
			if hook.Kind == "hijack" {
				return authResponse{checkResponse: true, status: status, err: err}
			}

			// Otherwise continue
			continue
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Hook (%s) in integration (%s) responded with action (%s)", hook.ID, hook.IntegrationID, res.Action), nil)

		switch res.Action {
		case config.ErrorHookResponse:
			// Simply log the error. No big deal. The hook can always hijack and throw error if this was serious
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Hook (%s) in integration (%s) sent error response - %s", hook.ID, hook.IntegrationID, res.Error), err, nil)

		case config.IgnoreHookResponse:
			// Do nothing
			break

		case config.HijackHookResponse:
			// Check if hook hook has permission to hijack
			if hook.Kind != "hijack" {
				helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Hook (%s) in integration (%s) does not has permission to hijack", hook.ID, hook.IntegrationID), nil)
				break
			}

			// We will skip hijack responses if a previous hook has already claimed this
			if hookResponse.checkResponse {
				helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot process highjack action of hook (%s) in integration (%s) since it has already been claimed by integration (%s)", hook.ID, hook.IntegrationID, hookResponse.integration), nil)
				break
			}

			// Set the hook Response
			hookResponse = authResponse{
				integration:   hook.IntegrationID,
				hook:          hook.ID,
				checkResponse: true,
				result:        res.Result,
				status:        status,
			}

			// Set error if exists
			if res.Error != "" {
				hookResponse.err = errors.New(res.Error)
			}
		}
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Exiting invoke hooks", map[string]interface{}{"resource": params.Resource, "verb": params.Op})

	return hookResponse
}

// HasPermissionForHook checks if a hook has the permission to hijack call
func HasPermissionForHook(config *config.IntegrationConfig, hook *config.IntegrationHook) bool {
	return hasPermissionForHook(config.ConfigPermissions, hook) || hasPermissionForHook(config.APIPermissions, hook)
}

func hasPermissionForHook(permissions []config.IntegrationPermission, hook *config.IntegrationHook) bool {
	for _, permission := range permissions {
		// Check if the resource types match
		if !utils.StringExists(permission.Resources, append(hook.Resource, "*")...) {
			continue
		}

		// Check if the attr match
		for k, v := range permission.Attributes {
			val, p := hook.Attributes[k]
			if !p {
				continue
			}

			if !utils.StringExists(v, append(val, "*")...) {
				continue
			}
		}

		// Check if the op matches
		if utils.StringExists(permission.Verbs, "*", hook.Kind) {
			return true
		}
	}

	return false
}
