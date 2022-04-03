package project

import (
	"context"

	"github.com/spacecloud-io/space-cloud/model"
)

// Hook implements the configman hook functionality
func (l *App) Hook(ctx context.Context, obj *model.ResourceObject) error {
	// TODO: allow only one aes key per project
	return nil
}
