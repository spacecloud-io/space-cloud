package file

import (
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (f *File) List(gvr schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}
	config, err := f.loadConfiguration()
	if err != nil {
		return list, err
	}

	key := source.GetModuleName(gvr)
	sources, ok := config[key]
	if !ok {
		return list, nil
	}

	list = common.ConvertToList(sources)
	return list, nil
}

func (f *File) Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	config, err := f.loadConfiguration()
	if err != nil {
		return nil, err
	}

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
