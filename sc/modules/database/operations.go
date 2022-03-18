package database

import (
	"context"
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman"
)

// Hook implements the configman hook functionality
func (l *App) Hook(ctx context.Context, obj *configman.ResourceObject) error {
	// Check if the resource belongs to this app
	if obj.Meta.Module != "database" {
		return fmt.Errorf("hook invoked for invalid resource type '%s/%s'", obj.Meta.Module, obj.Meta.Type)
	}

	// Process hook based on the resource type
	switch obj.Meta.Type {
	case "config":
		return processConfig(obj)
	case "schema":
		return l.processDBSchemaHook(ctx, obj)
	case "prepared-query":
		return processPreparedQuery(obj)
	default:
		return fmt.Errorf("hook invoked for invalid resource type '%s/%s'", obj.Meta.Module, obj.Meta.Type)
	}
}
