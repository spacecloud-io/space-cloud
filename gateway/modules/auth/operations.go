package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Encrypt encrypts a value if the aes key present in the config. The result is base64 encoded
// before being returned.
func (m *Module) Encrypt(value string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	return utils.Encrypt(m.aesKey, value)
}

// ParseToken simply parses and returns the claims of a provided token
func (m *Module) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	return m.jwt.ParseToken(ctx, token)
}

// GetAESKey gets aes key
func (m *Module) GetAESKey() []byte {
	m.RLock()
	defer m.RUnlock()
	return m.aesKey
}
