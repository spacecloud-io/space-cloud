package auth

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(App{})
	caddy.RegisterModule(AuthVerifyHandler{})
	caddy.RegisterModule(AuthOPAHandler{})
	caddy.RegisterModule(DenyUser{})
	caddy.RegisterModule(AuthenticateUser{})
}
