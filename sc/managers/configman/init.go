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
	caddy.RegisterModule(OperationHandler{})

	// All apps
	caddy.RegisterModule(ConfigMan{})
	caddy.RegisterModule(Store{})

}
