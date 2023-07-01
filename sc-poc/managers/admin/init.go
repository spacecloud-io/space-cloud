package admin

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(App{})
	caddy.RegisterModule(LoginHandler{})
	caddy.RegisterModule(RefreshHandler{})
	caddy.RegisterModule(AuthenticateSCUserPlugin{})
}
