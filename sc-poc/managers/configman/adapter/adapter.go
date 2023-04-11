package adapter

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
)

type Adapter interface {
	// Run starts the watcher.
	Run(context.Context) (chan common.ConfigType, error)

	// GetRawConfig returns the config in bytes.
	GetRawConfig() (common.ConfigType, error)
}
