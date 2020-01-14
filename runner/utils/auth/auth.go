package auth

import (
	"crypto/rsa"
	"errors"
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
	// JWT related stuff
	JWTAlgorithm JWTAlgorithm
	PublicKey    *rsa.PublicKey // for RSA
	Secret       string         // for HSA

	// For proxy authentication
	ProxySecret string

	Mode OperatingMode
}

// JWTAlgorithm describes the jwt algorithm to use
type JWTAlgorithm string

const (
	// RSA256 is used for rsa256 algorithm
	RSA256 JWTAlgorithm = "rsa256"

	// HS256 is used for hs256 algorithm
	HS256 JWTAlgorithm = "hs256"
)

// OperatingMode indicates the mode of operation
type OperatingMode string

const (
	// Runner indicates that the operating mode is runner
	Runner OperatingMode = "runner"

	// Server indicates that the operating mode is server
	Server OperatingMode = "server"
)

// New creates a new instance of the auth module
func New(config *Config) (*Module, error) {
	m := &Module{config: config}

	// The runner needs to fetch the public key from the server for rsa
	if config.JWTAlgorithm == RSA256 && config.Mode == "runner" {
		// Attempt fetching public key
		if success := m.fetchPublicKey(); !success {
			return nil, errors.New("could not initialise the auth module")
		}

		// Start the public key fetch routine
		go m.routineGetPublicKey()
	}

	return m, nil
}
