package auth

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/provider"
)

func init() {
	provider.Register("auth", 99)
	caddy.RegisterModule(App{})
	caddy.RegisterModule(AuthVerifyHandler{})
	caddy.RegisterModule(OPAPlugin{})
	caddy.RegisterModule(DenyUser{})
	caddy.RegisterModule(AuthenticateUser{})
}
