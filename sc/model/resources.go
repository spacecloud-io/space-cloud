package model

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/xeipuuv/gojsonschema"
)

type (
	// SCConfig is use to store sc config in a file
	SCConfig struct {
		Config map[string]ConfigModule `json:"config" yaml:"config" mapstructure:"config"`
	}

	// ConfigModule store module wise sc config
	ConfigModule map[string][]*ResourceObject

	// StoreMan implemments store
	StoreMan interface {
		ApplyResource(ctx context.Context, resourceObj *ResourceObject) error
		GetResource(ctx context.Context, meta *ResourceMeta) (*ResourceObject, error)
		GetResources(ctx context.Context, meta *ResourceMeta) (*ListResourceObjects, error)
		DeleteResource(ctx context.Context, meta *ResourceMeta) error
		DeleteResources(ctx context.Context, meta *ResourceMeta) error
	}

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

	// OperationTypes describes all the types which belong to a particular module
	OperationTypes map[string]*OperationTypeDefinition

	// OperationTypeDefinition describes the definition of a particular operation type
	OperationTypeDefinition struct {
		IsProtected     bool           `json:"isProtected" yaml:"isProtected"`
		Method          string         `json:"isSecure" yaml:"isSecure"`
		RequestSchema   interface{}    `json:"requestSchema" yaml:"requestSchema"`
		ResponseSchema  interface{}    `json:"responseSchema" yaml:"responseSchema"`
		RequiredParents []string       `json:"requiredParents" yaml:"requiredParents"`
		Controller      OperationHooks `json:"-" yaml:"-"`
	}

	// ConfigTypes describes all the types which belong to a particular module
	ConfigTypes map[string]*ConfigTypeDefinition

	// ConfigTypeDefinition describes the definition of a particular resource type
	ConfigTypeDefinition struct {
		IsSecure        bool        `json:"isSecure" yaml:"isSecure"`
		Schema          interface{} `json:"schema" yaml:"schema"`
		RequiredParents []string    `json:"requiredParents" yaml:"requiredParents"`
		Controller      ConfigHooks `json:"-" yaml:"-"`
	}

	// OperationHooks is used for processing an operation
	OperationHooks struct {
		DecodePayload func(ctx context.Context, reader io.ReadCloser) (interface{}, error)
		Handle        func(ctx context.Context, obj *ResourceObject, reqParams *RequestParams) (status int, payload interface{}, err error)
	}

	// ConfigHooks is used for processing config hooks
	ConfigHooks struct {
		PreApply   func(ctx context.Context, obj *ResourceObject, store StoreMan) error
		PostApply  func(ctx context.Context, obj *ResourceObject, store StoreMan) error
		PreGet     func(ctx context.Context, obj ResourceMeta, store StoreMan) error
		PostGet    func(ctx context.Context, obj *ListResourceObjects, store StoreMan) error
		PreDelete  func(ctx context.Context, obj ResourceMeta, store StoreMan) error
		PostDelete func(ctx context.Context, obj ResourceMeta, store StoreMan) error
	}

	// OperationCtrl are modules which expose operations
	OperationCtrl interface {
		GetOperationTypes() OperationTypes
	}

	// ConfigCtrl are modules which expose configs
	ConfigCtrl interface {
		GetConfigTypes() ConfigTypes
	}
)

// const (
// 	// PhasePreApply is invoked before the configuration is put in the store.
// 	PhasePreApply HookPhase = "pre-apply"

// 	// PhasePostApply is invoked after the configuration is put in the store.
// 	PhasePostApply HookPhase = "post-apply"

// 	// PhasePreGet is invoked before the configuration is get from the store.
// 	PhasePreGet HookPhase = "pre-get"

// 	// PhasePostGet is invoked after the configuration is get from the store.
// 	PhasePostGet HookPhase = "post-get"

// 	// PhasePreDelete is invoked before the configuration is delete from the store.
// 	PhasePreDelete HookPhase = "pre-apply"

// 	// PhasePostDelete is invoked after the configuration is delete from the store.
// 	PhasePostDelete HookPhase = "post-apply"
// )

// VerifyObject verifies if the config object is valid
func (typeDef *OperationTypeDefinition) VerifyObject(configObject *ResourceObject) ([]string, error) {
	// Check if all required fields are present
	if err := verifyMeta(configObject.Meta); err != nil {
		return nil, err
	}
	// Check if all required parents are present in object
	if err := verifyConfigParents(typeDef.RequiredParents, configObject.Meta.Parents); err != nil {
		return nil, err
	}

	// Very specification schema
	return verifySpecSchema(typeDef.RequestSchema, configObject.Spec)
}

// VerifyObject verifies if the config object is valid
func (typeDef *ConfigTypeDefinition) VerifyObject(configObject *ResourceObject) ([]string, error) {
	// Check if all required fields are present
	if err := verifyMeta(configObject.Meta); err != nil {
		return nil, err
	}
	// Check if all required parents are present in object
	if err := verifyConfigParents(typeDef.RequiredParents, configObject.Meta.Parents); err != nil {
		return nil, err
	}

	// Very specification schema
	return verifySpecSchema(typeDef.Schema, configObject.Spec)
}

func verifyMeta(meta ResourceMeta) error {
	if meta.Name == "" {
		return errors.New("resource name is missing")
	}
	if meta.Module == "" {
		return errors.New("resource module is missing")
	}
	if meta.Type == "" {
		return errors.New("resource type is missing")
	}
	return nil
}

func verifySpecSchema(schema, spec interface{}) ([]string, error) {
	// Skip verification if no json schema is supplied
	if schema == nil {
		return nil, nil
	}

	// Perform JSON schema validation
	schemaLoader := gojsonschema.NewGoLoader(schema)
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

func verifyConfigParents(requiredParents []string, parents map[string]string) error {
	// Simply return if object has no required parents
	if len(requiredParents) == 0 {
		return nil
	}

	// Send error if no parents are provided
	if parents == nil {
		return errors.New("resource doesn't have required parents")
	}

	// Check if all required parents are available
	for _, parent := range requiredParents {
		if _, p := parents[parent]; !p {
			return fmt.Errorf("parent '%s' not present in resource", parent)
		}
	}

	return nil
}
