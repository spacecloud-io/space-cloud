package letsencrypt

import (
	"fmt"

	"github.com/mholt/certmagic"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// LetsEncrypt manages letsencrypt certificates
type LetsEncrypt struct {

	// For internal use
	client *certmagic.Config
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
	case config.LetsEncryptStoreLocal:
		client.Storage = certmagic.Default.Storage
	case config.LetsEncryptStoreSC:
		// TODO: implement this
	default:
		return nil, fmt.Errorf("unsupported store type (%s) provided for lets encrypt", c.StoreType)
	}

	// Clean whitelisted domains
	if c.WhitelistedDomains == nil {
		c.WhitelistedDomains = make([]string, 0)
	}

	if err := client.ManageSync(c.WhitelistedDomains); err != nil {
		return nil, err
	}

	return &LetsEncrypt{client: client}, nil
}
