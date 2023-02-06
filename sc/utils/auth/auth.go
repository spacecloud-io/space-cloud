package auth

import "github.com/spacecloud-io/space-cloud/config"

// Module handles everything related to JWT Tokens
type Module struct {
	secrets []*config.Secret
}

// New creates a new module instance
func New(secrets []*config.Secret) *Module {
	return &Module{secrets: secrets}
}
