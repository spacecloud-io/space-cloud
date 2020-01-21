package letsencrypt

import (
	"fmt"
	"sync"

	"github.com/mholt/certmagic"
)

// LetsEncrypt manages letsencrypt certificates
type LetsEncrypt struct {
	lock sync.Mutex

	// For internal use
	client  *certmagic.Config
	domains domainMapping
}

// New creates a new letsencrypt module
func New() (*LetsEncrypt, error) {
	// Load config from environment variables
	c := loadConfig()

	client := certmagic.NewDefault()
	client.Agreed = true
	client.Email = c.Email

	// Set the store for certificates
	switch c.StoreType {
	case StoreLocal:
		client.Storage = certmagic.Default.Storage
	case StoreSC:
		client.Storage = NewScStore()
	default:
		return nil, fmt.Errorf("unsupported store type (%s) provided for lets encrypt", c.StoreType)
	}

	return &LetsEncrypt{client: client, domains: domainMapping{}}, nil
}
