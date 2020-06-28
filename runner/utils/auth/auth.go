package auth

import (
	"sync"
)

// Module manages the auth module
type Module struct {
	lock sync.RWMutex

	// For internal use
	config *Config
}

// Config is the object used to configure the auth module
type Config struct {
	// For authentication of runner and store
	Secret string

	// For proxy authentication
	ProxySecret string

	// disable authentication while development
	IsDev bool
}

// JWTAlgorithm describes the jwt algorithm to use
type JWTAlgorithm string

const (
	// RSA256 is used for rsa256 algorithm
	RSA256 JWTAlgorithm = "rsa256"

	// HS256 is used for hs256 algorithm
	HS256 JWTAlgorithm = "hs256"
)

// New creates a new instance of the auth module
func New(config *Config) (*Module, error) {
	m := &Module{config: config}

	return m, nil
}
