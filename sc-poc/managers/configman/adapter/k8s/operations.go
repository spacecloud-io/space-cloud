package k8s

import (
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all the registered sources of a particular source type
func (k *K8s) List(gvr schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}
	key := source.GetModuleName(gvr)

	config := k.getConfig()
	sources, ok := config[key]
	if !ok {
		return list, nil
	}

	list = common.ConvertToList(sources)
	return list, nil
}

func (k *K8s) Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	key := source.GetModuleName(gvr)

	config := k.getConfig()
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
