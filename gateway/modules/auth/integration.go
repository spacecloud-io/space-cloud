package auth

import (
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetIntegrationToken returns a token for the integration module
func (m *Module) GetIntegrationToken(id string) (string, error) {
	return m.CreateToken(map[string]interface{}{"id": id, "role": "integration"})
}

// GetMissionControlToken returns a token to be used by mission control
func (m *Module) GetMissionControlToken(claims map[string]interface{}) (string, error) {
	return m.CreateToken(map[string]interface{}{"id": utils.InternalUserID, "claims": claims})
}
