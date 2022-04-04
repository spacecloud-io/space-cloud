package file

import (
	"context"
	"fmt"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

// File implements file store
type File struct {
	Path string `json:"path,omitempty"`

	logger *zap.Logger
}

type config struct {
	Config map[string]configModule `json:"config" yaml:"config" mapstructure:"config"`
}

type configModule map[string][]*model.ResourceObject

// ApplyResource applies resource in the store
func (f File) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.Path, scConfig); err != nil {
		return err
	}

	if scConfig.Config == nil {
		scConfig.Config = make(map[string]configModule)
	}

	module, ok := scConfig.Config[resourceObj.Meta.Module]
	if !ok {
		module = make(configModule)
	}

	moduleType, ok := module[resourceObj.Meta.Type]
	if !ok {
		moduleType = make([]*model.ResourceObject, 0)
	}

	moduleType = addResource(moduleType, resourceObj)
	module[resourceObj.Meta.Type] = moduleType
	scConfig.Config[resourceObj.Meta.Module] = module

	return utils.StoreFile(f.Path, scConfig)
}

func addResource(moduleType []*model.ResourceObject, resourceObj *model.ResourceObject) []*model.ResourceObject {
	for index, resource := range moduleType {
		if resource.Meta.Name == resourceObj.Meta.Name && matchParent(resource.Meta.Parents, resourceObj.Meta.Parents) {
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

	moduleType, ok := module[meta.Type]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	for _, resourceObj := range moduleType {
		if meta.Name == resourceObj.Meta.Name && matchParent(meta.Parents, resourceObj.Meta.Parents) {
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

	moduleType, ok := module[meta.Type]
	if !ok {
		return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	resourceList := make([]*model.ResourceObject, 0)
	for _, resourceObj := range moduleType {
		if matchParent(meta.Parents, resourceObj.Meta.Parents) {
			resourceList = append(resourceList, resourceObj)
		}
	}

	return &model.ListResourceObjects{List: resourceList}, nil
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

	moduleType, ok := module[meta.Type]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	for i, resourceObj := range moduleType {
		if meta.Name == resourceObj.Meta.Name && matchParent(meta.Parents, resourceObj.Meta.Parents) {
			moduleType = append(moduleType[:i], moduleType[i+1:]...)
			module[meta.Type] = moduleType
			scConfig.Config[meta.Module] = module
			break
		}
	}
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

	moduleType, ok := module[meta.Type]
	if !ok {
		return fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
	}

	moduleType = deleteResources(moduleType, meta)
	module[meta.Type] = moduleType
	scConfig.Config[meta.Module] = module

	return utils.StoreFile(f.Path, scConfig)
}

func (f File) SetLogger(logger *zap.Logger) {
	f.logger = logger
}

func matchParent(a, b map[string]string) bool {
	for k, v := range b {
		v1, ok := a[k]
		if !ok || v1 != v {
			return false
		}
	}

	return true
}

func deleteResources(moduleType []*model.ResourceObject, meta *model.ResourceMeta) []*model.ResourceObject {
	tempModuleType := make([]*model.ResourceObject, 0)

	for _, resourceObj := range moduleType {
		if !matchParent(resourceObj.Meta.Parents, meta.Parents) {
			tempModuleType = append(tempModuleType, resourceObj)
		}
	}
	return tempModuleType
}
