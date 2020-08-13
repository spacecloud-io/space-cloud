package integration

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func isIntegrationRequest(claims map[string]interface{}) bool {
	role, p := claims["role"]

	// The role must be present in the claims
	if !p {
		return false
	}

	// The role must be integration
	if role != "integration" {
		return false
	}

	return true
}

func (m *Manager) checkPermissions(kind, resource, op string, claims map[string]interface{}, attr map[string]string) error {
	// Extract the necessary claims
	id := claims["id"]

	i, p := m.config[id.(string)]
	if !p {
		return utils.LogError(fmt.Sprintf("Integration (%s) not found", id), module, checkAuth, nil)
	}

	// Get the write permission block
	var permissions []config.IntegrationPermission
	switch kind {
	case "config":
		permissions = i.ConfigPermissions
	case "api":
		permissions = i.APIPermissions
	default:
		return utils.LogError(fmt.Sprintf("Invalid permission type (%s) provided", kind), module, checkAuth, nil)
	}

	// Check if the integration has the required permission

	for _, permission := range permissions {
		// Check if the resource types match
		if !utils.StringExists(permission.Resources, "*", resource) {
			continue
		}

		// Check if the op matches
		if !utils.StringExists(permission.Verbs, "*", op) {
			continue
		}

		// Return if attr is nil since all other things matched
		if attr == nil {
			return nil
		}

		// Check if the attr match
		for k, v := range permission.Attributes {
			val, p := attr[k]
			if !p {
				continue
			}

			if !utils.StringExists(v, "*", val) {
				continue
			}
		}

		return nil
	}

	return utils.LogError(fmt.Sprintf("Integration (%s) does not have the required permissions", id), module, checkAuth, nil)
}
