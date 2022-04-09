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
	path string `json:"path,omitempty"`

	logger *zap.Logger
}

func New(logger *zap.Logger, path string) (*File, error) {
	return &File{logger: logger, path: path}, nil
}

// ApplyResource applies resource in the store
func (f *File) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.path, scConfig); err != nil {
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

	return utils.StoreFile(f.path, scConfig)
}

// GetResource gets resource from the store
func (f *File) GetResource(ctx context.Context, meta *model.ResourceMeta) (*model.ResourceObject, error) {
	scConfig := new(config)
	if err := utils.LoadFile(f.path, scConfig); err != nil {
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
		if meta.Name == resourceObj.Meta.Name && matchParent(resourceObj.Meta.Parents, meta.Parents) {
			return resourceObj, nil
		}
	}

	return nil, fmt.Errorf("no resource found for %s - %s - %s", meta.Module, meta.Type, meta.Name)
}

// GetResources gets resources from the store
func (f *File) GetResources(ctx context.Context, meta *model.ResourceMeta) (*model.ListResourceObjects, error) {
	scConfig := new(config)
	if err := utils.LoadFile(f.path, scConfig); err != nil {
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
		if matchParent(resourceObj.Meta.Parents, meta.Parents) {
			resourceList = append(resourceList, resourceObj)
		}
	}

	return &model.ListResourceObjects{List: resourceList}, nil
}

// DeleteResource delete resource from the store
func (f *File) DeleteResource(ctx context.Context, meta *model.ResourceMeta) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.path, scConfig); err != nil {
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
		if meta.Name == resourceObj.Meta.Name && matchParent(resourceObj.Meta.Parents, meta.Parents) {
			moduleType = append(moduleType[:i], moduleType[i+1:]...)
			module[meta.Type] = moduleType
			scConfig.Config[meta.Module] = module
			break
		}
	}
	return utils.StoreFile(f.path, scConfig)
}

// DeleteResources delete resources from the store
func (f *File) DeleteResources(ctx context.Context, meta *model.ResourceMeta) error {
	scConfig := new(config)
	if err := utils.LoadFile(f.path, scConfig); err != nil {
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

	return utils.StoreFile(f.path, scConfig)
}

func (f *File) SetLogger(logger *zap.Logger) {
	f.logger = logger
}

func (f *File) Destruct() error {
	return nil
}
