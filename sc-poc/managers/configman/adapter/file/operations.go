package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
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

	key := source.GetModuleName(gvr)
	config := f.getConfig()
	newArr := []*unstructured.Unstructured{}

	for _, spec := range config[key] {
		if spec.GetName() != name {
			newArr = append(newArr, spec)
		}
	}

	fileName := generateYAMLFileName(gvr.Group, gvr.Resource)
	// If length of array is 0, delete the file.
	if len(newArr) == 0 {
		if err := os.Remove(filepath.Join(f.path, fileName)); err != nil {
			return fmt.Errorf("could not delete file: %v", err)
		}

		return nil
	}

	// Overwrite the file with the new config
	var newSpec []byte
	for _, spec := range newArr {
		// Get spec in bytes
		data, err := utils.GetBytesFromSpec(spec)
		if err != nil {
			return err
		}

		newSpec = append(newSpec, data...)
	}

	err := os.WriteFile(filepath.Join(f.path, fileName), newSpec, 0777)
	if err != nil {
		return err
	}

	return nil
}
