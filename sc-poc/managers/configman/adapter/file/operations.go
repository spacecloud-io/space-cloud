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
	list := &unstructured.UnstructuredList{}
	config := f.getConfig()

	key := source.GetModuleName(gvr)
	sources, ok := config[key]
	if !ok {
		return list, nil
	}

	list = common.ConvertToList(sources)
	return list, nil
}

// Get returns a registered source
func (f *File) Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	config := f.getConfig()

	key := source.GetModuleName(gvr)
	sources, ok := config[key]
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

	fileName := gvr.Resource + ".yaml"
	// Check if the source file exists
	_, err = os.Stat(filepath.Join(f.Path, fileName))
	// If source file does not exists create a new one and
	// write this spec into this file
	if os.IsNotExist(err) {
		err := os.WriteFile(filepath.Join(f.Path, fileName), data, 0777)
		if err != nil {
			return err
		}
	} else {
		// If source file exists append the spec.
		err := utils.AppendToFile(filepath.Join(f.Path, fileName), data)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes a source
func (f *File) Delete(gvr schema.GroupVersionResource, name string) error {
	if err := f.reorganizeFileStructure(); err != nil {
		return err
	}

	fileName := gvr.Resource + ".yaml"
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
