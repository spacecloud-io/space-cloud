package configman

import (
	"context"
	"fmt"
	"strings"

	"github.com/spacecloud-io/space-cloud/model"
)

func applyHooks(ctx context.Context, module string, typeDef *model.TypeDefinition, phase model.HookPhase, loadApp loadApp, resourceObj *model.ResourceObject) error {
	// Invoke hooks if any
	hook, err := loadHook(module, typeDef, model.PhasePostApply, loadApp)
	if err != nil {
		return err
	}

	// Invoke hook if exists
	if hook != nil {
		if err := hook.Hook(ctx, resourceObj); err != nil {
			return err
		}
	}
	return nil
}

func extractPathParams(urlPath string) (op, module, typeName, resourceName string, err error) {
	// Set the default operation to single
	op = "single"

	// Check if url has proper length
	arr := strings.Split(urlPath[1:], "/")
	if len(arr) > 5 || len(arr) < 4 {
		err = fmt.Errorf("invalid config url provided - %s", urlPath)
		return
	}

	// Check the operation type
	if len(arr) == 5 {
		op = "list"
		resourceName = arr[4]
	}

	// Set the other parameters
	module = arr[2]
	typeName = arr[3]
	return
}

func prepareErrorResponseBody(err error, schemaErrors []string) interface{} {
	return map[string]interface{}{
		"error":        err.Error(),
		"schemaErrors": schemaErrors,
	}
}
