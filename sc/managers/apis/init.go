package apis

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(RootAPIHandler{})
	caddy.RegisterModule(APIHandler{})
	caddy.RegisterModule(CorsHandler{})
}
