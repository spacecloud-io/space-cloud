package model

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
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

	// ListResourceObjects describes the configuration for resources
	ListResourceObjects struct {
		List []*ResourceObject `json:"list" yaml:"list"`
	}

	// Types describes all the types which belong to a particular module
	Types map[string]*TypeDefinition

	// TypeDefinition describes the definition of a particular resource type
	TypeDefinition struct {
		IsSecure        bool        `json:"isSecure" yaml:"isSecure"`
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

	// PhasePostApply is invoked after the configuration is put in the store.
	PhasePostApply HookPhase = "post-apply"

	// PhasePreGet is invoked before the configuration is get from the store.
	PhasePreGet HookPhase = "pre-get"

	// PhasePostGet is invoked after the configuration is get from the store.
	PhasePostGet HookPhase = "post-get"

	// PhasePreDelete is invoked before the configuration is delete from the store.
	PhasePreDelete HookPhase = "pre-apply"

	// PhasePostDelete is invoked after the configuration is delete from the store.
	PhasePostDelete HookPhase = "post-apply"
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

func verifySpecSchema(typeDef *TypeDefinition, spec interface{}) ([]string, error) {
	// Skip verification if no json schema is supplied
	if typeDef.Schema == nil {
		return nil, nil
	}

	// Perform JSON schema validation
	schemaLoader := gojsonschema.NewGoLoader(typeDef.Schema)
	documentLoader := gojsonschema.NewGoLoader(spec)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	// Skip if no errros were found
	if result.Valid() {
		return nil, nil
	}

	// Send back all the errors
	arr := make([]string, len(result.Errors()))
	for i, err := range result.Errors() {
		arr[i] = err.String()
	}

	return arr, fmt.Errorf("json schema validation failed")
}

func verifyConfigParents(typeDef *TypeDefinition, parents map[string]string) error {
	// Simply return if object has no required parents
	if len(typeDef.RequiredParents) == 0 {
		return nil
	}

	// Send error if no parents are provided
	if parents == nil {
		return errors.New("resource doesn't have required parents")
	}

	// Check if all required parents are available
	for _, parent := range typeDef.RequiredParents {
		if _, p := parents[parent]; !p {
			return fmt.Errorf("parent '%s' not present in resource", parent)
		}
	}

	return nil
}
