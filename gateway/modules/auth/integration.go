package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetIntegrationToken returns a token for the integration module
func (m *Module) GetIntegrationToken(ctx context.Context, id string) (string, error) {
	return m.CreateToken(ctx, map[string]interface{}{"id": id, "role": "integration"})
}

// GetMissionControlToken returns a token to be used by mission control
func (m *Module) GetMissionControlToken(ctx context.Context, claims map[string]interface{}) (string, error) {
	return m.CreateToken(context.Background(), map[string]interface{}{"id": utils.InternalUserID, "claims": claims})
}
