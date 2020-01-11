package letsencrypt

import (
	"os"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetConfig sets the config required by lets encrypt
func (l *LetsEncrypt) SetConfig(whitelist []string) error {
	if whitelist == nil {
		whitelist = make([]string, 0)
	}

	return l.client.ManageSync(whitelist)
}

func loadConfig() *config.LetsEncrypt {
	prefix := "LETSENCRYPT_"
	c := &config.LetsEncrypt{WhitelistedDomains: []string{}, StoreType: config.LetsEncryptStoreLocal}
	if email, p := os.LookupEnv(prefix + "EMAIL"); p {
		c.Email = email
	}

	if whitelist, p := os.LookupEnv(prefix + "WHITELIST"); p {
		c.WhitelistedDomains = strings.Split(whitelist, ",")
	}

	if store, p := os.LookupEnv(prefix + "STORE"); p {
		c.StoreType = config.LetsEncryptStoreType(store)
	}

	return c
}
