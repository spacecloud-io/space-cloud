package configman

import (
	"context"
	"errors"
)

type (
	loadApp func(appName string) (interface{}, error)

	// HookImpl is a controller which implements the hook interface
	HookImpl interface {
		Hook(ctx context.Context, obj *ResourceObject) error
	}
)

type (
	// ResourceMeta contains the meta information about a resource
	ResourceMeta struct {
		Module  string            `json:"module" yaml:"module"`
		Type    string            `json:"type" yaml:"type"`
		Name    string            `json:"name" yaml:"name"`
		Parents map[string]string `json:"parent" yaml:"parent"`
	}

	// ResourceObject describes the configuration for a single resource
	ResourceObject struct {
		Meta ResourceMeta `json:"meta" yaml:"meta"`
		Spec interface{}  `json:"spec" yaml:"spec"`
	}

	// Types describes all the types which belong to a particular module
	Types map[string]*TypeDefinition

	// TypeDefinition describes the definition of a particular resource type
	TypeDefinition struct {
		Schema          interface{} `json:"schema" yaml:"schema"`
		Hooks           Hooks       `json:"hooks" yaml:"hooks"`
		RequiredParents []string    `json:"requiredParents" yaml:"requiredParents"`
	}

	// Hooks describes the hooks to be invoked on the module
	Hooks map[HookPhase]struct{}

	// HookPhase describes the various phases where hooks can be invoked.
	HookPhase string
)

const (
	// PhasePreApply is invoked before the configuration is put in the store.
	PhasePreApply HookPhase = "pre-apply"
)

// VerifyObject verifies if the config object is valid
func (typeDef *TypeDefinition) VerifyObject(configObject *ResourceObject) ([]string, error) {
	// Check if all required fields are present
	if configObject.Meta.Name == "" {
		return nil, errors.New("resource name is missing")
	}
	if configObject.Meta.Module == "" {
		return nil, errors.New("resource module is missing")
	}
	if configObject.Meta.Type == "" {
		return nil, errors.New("resource type is missing")
	}

	// Check if all required parents are present in object
	if err := verifyConfigParents(typeDef, configObject.Meta.Parents); err != nil {
		return nil, err
	}

	// Very specification schema
	return verifySpecSchema(typeDef, configObject.Spec)
}
