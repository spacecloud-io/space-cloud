package handlers

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(List{})
}
