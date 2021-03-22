package letsencrypt

import (
	"context"
	"fmt"
	"sync"

	"github.com/caddyserver/certmagic"
	"github.com/spaceuptech/helpers"
)

// LetsEncrypt manages letsencrypt certificates
type LetsEncrypt struct {
	lock sync.Mutex

	// For internal use
	config  *certmagic.Config
	client  *certmagic.ACMEManager
	domains domainMapping
}

// New creates a new letsencrypt module
func New() (*LetsEncrypt, error) {
	// Load config from environment variables
	c := loadConfig()

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = c.Email

	config := certmagic.NewDefault()

	// Set the store for certificates
	switch c.StoreType {
	case StoreLocal:
		config.Storage = certmagic.Default.Storage
	case StoreSC:
		config.Storage = NewScStore()
	case StoreKube:
		c, err := NewKubeStore()
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to initialize kubernetes store for let encrypt", err, nil)
		}
		config.Storage = c
	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unsupported store type (%s) provided for lets encrypt", c.StoreType), nil, nil)
	}

	client := certmagic.NewACMEManager(config, certmagic.ACMEManager{
		Agreed: true,
	})

	return &LetsEncrypt{config: config, client: client, domains: domainMapping{}}, nil
}

// SetLetsEncryptEmail sets config email
func (l *LetsEncrypt) SetLetsEncryptEmail(email string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.client.Email = email
}
