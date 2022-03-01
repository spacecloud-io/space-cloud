package configman

import "github.com/caddyserver/caddy/v2"

func init() {
	// config modules
	caddy.RegisterModule(ConfigGetHandler{})
	caddy.RegisterModule(ConfigDeleteHandler{})
	caddy.RegisterModule(ConfigPostHandler{})

	// operation modules
	caddy.RegisterModule(OperationGetHandler{})
	caddy.RegisterModule(OperationDeleteHandler{})
	caddy.RegisterModule(OperationPostHandler{})

}
