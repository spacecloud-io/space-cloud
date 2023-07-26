package auth

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(App{})
	caddy.RegisterModule(JWTAuthVerifyHandler{})
	caddy.RegisterModule(KratosAuthVerifyHandler{})
	caddy.RegisterModule(AuthOPAHandler{})
}
