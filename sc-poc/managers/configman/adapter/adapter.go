package adapter

import "context"

type Adapter interface {
	// Run starts the watcher.
	Run(context.Context) (chan []byte, error)

	// GetRawConfig returns the config in bytes.
	GetRawConfig() ([]byte, error)
}
