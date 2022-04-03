package configman

import "github.com/caddyserver/caddy/v2"

func init() {
	// openapi modules
	caddy.RegisterModule(OpenAPIHandler{})

	// config modules
	caddy.RegisterModule(ConfigGetHandler{})
	caddy.RegisterModule(ConfigDeleteHandler{})
	caddy.RegisterModule(ConfigApplyHandler{})

	// operation modules
	caddy.RegisterModule(OperationGetHandler{})
	caddy.RegisterModule(OperationDeleteHandler{})
	caddy.RegisterModule(OperationPostHandler{})

	// store module
	caddy.RegisterModule(ConfigMan{})

}
