package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all registered sources of a specific source type
func (f *File) List(gvr schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	list := &unstructured.UnstructuredList{}
	key := source.GetModuleName(gvr)
	sources, ok := f.configuration[key]
	if !ok {
		return list, nil
	}

	list = common.ConvertToList(sources)
	return list, nil
}

// Get returns a registered source
func (f *File) Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	key := source.GetModuleName(gvr)
	sources, ok := f.configuration[key]
	if !ok {
		return nil, fmt.Errorf("source with name: %s not found", name)
	}

	for _, src := range sources {
		if src.GetName() == name {
			return src, nil
		}
	}

	return nil, fmt.Errorf("source with name: %s not found", name)
}

// Apply creates/updates a source
func (f *File) Apply(gvr schema.GroupVersionResource, spec *unstructured.Unstructured) error {
	if err := f.reorganizeFileStructure(); err != nil {
		return err
	}

	// Get spec in bytes
	data, err := utils.GetBytesFromSpec(spec)
	if err != nil {
		return err
	}

	return f.persistConfig(gvr, data)
}

// Delete deletes a source
func (f *File) Delete(gvr schema.GroupVersionResource, name string) error {
	if err := f.reorganizeFileStructure(); err != nil {
		return err
	}

	fileName := generateYAMLFileName(gvr.Group, gvr.Resource)
	// Parse the config file
	arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(f.Path, fileName))
	if err != nil {
		f.logger.Error("Unable to parse config file", zap.String("file", fileName), zap.Error(err))
		return err
	}

	newArr := []*unstructured.Unstructured{}
	for _, spec := range arr {
		if spec.GetName() != name {
			newArr = append(newArr, spec)
		}
	}

	// Remove the old file
	if err = os.Remove(filepath.Join(f.Path, fileName)); err != nil {
		return fmt.Errorf("could not delete file: %v", err)
	}

	// Create new file and append the new specs
	if len(newArr) != 0 {
		if _, err := os.Create(filepath.Join(f.Path, fileName)); err != nil {
			return err
		}

		for _, spec := range newArr {
			// Get spec in bytes
			data, err := utils.GetBytesFromSpec(spec)
			if err != nil {
				return err
			}

			err = utils.AppendToFile(filepath.Join(f.Path, fileName), data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
