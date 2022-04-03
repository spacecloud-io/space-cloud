package file

import (
	"context"
	"fmt"
	"reflect"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

// File implements file store
type File struct {
	Path string `json:"path,omitempty"`
}

type config struct {
	Config map[string]*configModule `json:"config" yaml:"config" mapstructure:"config"`
}

type configModule struct {
	ModuleType map[string][]*model.ResourceObject `json:"type" yaml:"type" mapstructure:"type"`
}

// ApplyResource applies resource in the store
func (f File) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return err
	}

	if scConfig.Config == nil {
		scConfig.Config = make(map[string]*configModule)
	}

	module, ok := scConfig.Config[resourceObj.Meta.Module]
	if !ok {
		module.ModuleType = make(map[string][]*model.ResourceObject)
		module.ModuleType[resourceObj.Meta.Type] = make([]*model.ResourceObject, 0)
	}

	moduleType, ok := module.ModuleType[resourceObj.Meta.Type]
	if !ok {
		moduleType = make([]*model.ResourceObject, 0)
	}

	moduleType = addResource(moduleType, resourceObj)
	module.ModuleType[resourceObj.Meta.Type] = moduleType
	scConfig.Config[resourceObj.Meta.Module] = module

	return utils.StoreFile(f.Path, scConfig)
}

func addResource(moduleType []*model.ResourceObject, resourceObj *model.ResourceObject) []*model.ResourceObject {
	for index, resource := range moduleType {
		if resource.Meta.Name == resourceObj.Meta.Name && reflect.DeepEqual(resource.Meta.Parents, resourceObj.Meta.Parents) {
			moduleType[index] = resourceObj
			return moduleType
		}
	}

	moduleType = append(moduleType, resourceObj)
	return moduleType
}

// GetResource gets resource from the store
func (f File) GetResource(ctx context.Context, meta *model.ResourceMeta) (*model.ResourceObject, error) {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return nil, err
	}

	module, ok := scConfig.Config[meta.Module]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType, ok := module.ModuleType[meta.Type]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	for _, resourceObj := range moduleType {
		if meta.Name == resourceObj.Meta.Name && reflect.DeepEqual(meta.Parents, resourceObj.Meta.Parents) {
			return resourceObj, nil
		}
	}

	return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
}

// GetResources gets resources from the store
func (f File) GetResources(ctx context.Context, meta *model.ResourceMeta) (*model.ListResourceObjects, error) {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return nil, err
	}

	module, ok := scConfig.Config[meta.Module]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType, ok := module.ModuleType[meta.Type]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	return &model.ListResourceObjects{List: moduleType}, nil
}

// DeleteResource delete resource from the store
func (f File) DeleteResource(ctx context.Context, meta *model.ResourceMeta) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return err
	}

	module, ok := scConfig.Config[meta.Module]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType, ok := module.ModuleType[meta.Type]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	index := -1
	for i, resourceObj := range moduleType {
		if meta.Name == resourceObj.Meta.Name && reflect.DeepEqual(meta.Parents, resourceObj.Meta.Parents) {
			index = i
		}
	}

	if index == -1 {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType = append(moduleType[:index], moduleType[index+1:]...)
	module.ModuleType[meta.Type] = moduleType
	scConfig.Config[meta.Module] = module

	return utils.StoreFile(f.Path, scConfig)
}

// DeleteResources delete resources from the store
func (f File) DeleteResources(ctx context.Context, meta *model.ResourceMeta) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return err
	}

	module, ok := scConfig.Config[meta.Module]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType, ok := module.ModuleType[meta.Type]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	tempModuleType := deleteResources(moduleType, meta)
	if len(tempModuleType) == len(moduleType) {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}
	module.ModuleType[meta.Type] = tempModuleType
	scConfig.Config[meta.Module] = module

	return utils.StoreFile(f.Path, scConfig)
}

func deleteResources(moduleType []*model.ResourceObject, meta *model.ResourceMeta) []*model.ResourceObject {
	tempModuleType := make([]*model.ResourceObject, 0)

	for _, resourceObj := range moduleType {
		if !reflect.DeepEqual(meta.Parents, resourceObj.Meta.Parents) {
			tempModuleType = append(tempModuleType, resourceObj)
		}
	}
	return tempModuleType
}
