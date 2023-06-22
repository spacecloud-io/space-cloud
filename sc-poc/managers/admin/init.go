package admin

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(ListSources{})
	caddy.RegisterModule(Login{})
	caddy.RegisterModule(Refresh{})
}
