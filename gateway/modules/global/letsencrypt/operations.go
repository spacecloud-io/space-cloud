package letsencrypt

import (
	"crypto/tls"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// TLSConfig returns the tls config to be used by the http server
func (l *LetsEncrypt) TLSConfig() *tls.Config {
	return l.config.TLSConfig()
}

// LetsEncryptHTTPChallengeHandler handle the http challenge
func (l *LetsEncrypt) LetsEncryptHTTPChallengeHandler(h http.Handler) http.Handler {
	return l.client.HTTPChallengeHandler(h)
}

// AddExistingCertificate lets the user add an existing certificate. This certificate
// will not be automatically renewed via let's encrypt
func (l *LetsEncrypt) AddExistingCertificate(certFile, keyFile string) error {
	return l.config.CacheUnmanagedCertificatePEMFile(certFile, keyFile, []string{})
}

// SetProjectDomains sets the config required by lets encrypt
func (l *LetsEncrypt) SetProjectDomains(project string, c *config.LetsEncrypt) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// if c == nil {
	// 	c = new(config.LetsEncrypt)
	// }

	if c.WhitelistedDomains == nil {
		c.WhitelistedDomains = make([]string, 0)
	}

	if len(c.WhitelistedDomains) == 0 {
		return nil
	}

	l.domains.setProjectDomains(project, c.WhitelistedDomains)
	return l.config.ManageSync(l.domains.getUniqueDomains())
}

// DeleteProjectDomains deletes a projects associated domains
func (l *LetsEncrypt) DeleteProjectDomains(project string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.domains.deleteProject(project)
	return l.config.ManageSync(l.domains.getUniqueDomains())
}
