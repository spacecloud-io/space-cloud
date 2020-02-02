package letsencrypt

import (
	"os"
	"strings"
)

// Config describes the configuration for let's encrypt
type Config struct {
	Email              string
	WhitelistedDomains []string
	StoreType          StoreType
}

// StoreType describes the store used by the lets encrypt module
type StoreType string

const (
	// StoreLocal is used when the local filesystem us used as a store
	StoreLocal StoreType = "local"

	// StoreSC is used when SC is supposed to be used as a store
	StoreSC StoreType = "sc"

	// StoreKube is used when kubernetes is supposed to be used as a store
	StoreKube StoreType = "kube"
)

func loadConfig() *Config {
	prefix := "LETSENCRYPT_"
	c := &Config{WhitelistedDomains: []string{}, StoreType: StoreLocal}
	if email, p := os.LookupEnv(prefix + "EMAIL"); p {
		c.Email = email
	}

	if whitelist, p := os.LookupEnv(prefix + "WHITELIST"); p {
		c.WhitelistedDomains = strings.Split(whitelist, ",")
	}

	if store, p := os.LookupEnv(prefix + "STORE"); p {
		c.StoreType = StoreType(store)
	}

	return c
}
