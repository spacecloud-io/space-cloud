package auth

import (
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Encrypt encrypts a value if the aes key present in the config. The result is base64 encoded
// before being returned.
func (m *Module) Encrypt(value string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	return utils.Encrypt(m.aesKey, value)
}
