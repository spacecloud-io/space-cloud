package file

import (
	"github.com/spacecloud-io/space-cloud/managers/source"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (f *File) List(gvr schema.GroupVersionResource) ([]*unstructured.Unstructured, error) {
	config, err := f.loadConfiguration()
	if err != nil {
		return nil, err
	}

	key := source.GetModuleName(gvr)
	val, ok := config[key]
	if !ok {
		return []*unstructured.Unstructured{}, nil
	}

	return val, nil
}
