package letsencrypt

import (
	"crypto/tls"
	"net/http"
)

// TLSConfig returns the tls config to be used by the http server
func (l *LetsEncrypt) TLSConfig() *tls.Config {
	return l.client.TLSConfig()
}

// LetsEncryptHTTPChallengeHandler handle the http challenge
func (l *LetsEncrypt) LetsEncryptHTTPChallengeHandler(h http.Handler) http.Handler {
	return l.client.HTTPChallengeHandler(h)
}

func (l *LetsEncrypt) AddExistingCertificate(certFile, keyFile string) error {
	return l.client.CacheUnmanagedCertificatePEMFile(certFile, keyFile, []string{})
}
