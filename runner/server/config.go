package server

import (
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver"
)

// Config is the object required to configure the runner
type Config struct {
	Port             string
	ProxyPort        string
	IsMetricDisabled bool

	// Configuration for the driver
	Driver *driver.Config

	// Configuration for the auth module
	Auth *auth.Config
}
