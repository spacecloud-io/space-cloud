package letsencrypt

import (
	"fmt"
	"sync"

	"github.com/mholt/certmagic"
	"github.com/spaceuptech/helpers"
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
	case StoreKube:
		c, err := NewKubeStore()
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to initialize kubernetes store for let encrypt", err, nil)
		}
		client.Storage = c
	default:
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unsupported store type (%s) provided for lets encrypt", c.StoreType), nil, nil)
	}

	return &LetsEncrypt{client: client, domains: domainMapping{}}, nil
}

//SetLetsEncryptEmail sets config email
func (l *LetsEncrypt) SetLetsEncryptEmail(email string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.client.Email = email
}
