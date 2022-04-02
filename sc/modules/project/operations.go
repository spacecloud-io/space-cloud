package project

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/configman"
)

// Hook implements the configman hook functionality
func (l *App) Hook(ctx context.Context, obj *configman.ResourceObject) error {
	// TODO: allow only one aes key per project
	return nil
}
