package k8s

import (
	"github.com/spacecloud-io/space-cloud/managers/source"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all the registered sources of a particular source type
func (k *K8s) List(gvr schema.GroupVersionResource) ([]*unstructured.Unstructured, error) {
	k.lock.RLock()
	defer k.lock.RUnlock()

	key := source.GetModuleName(gvr)
	val, ok := k.configuration[key]
	if !ok {
		return []*unstructured.Unstructured{}, nil
	}

	return val, nil
}
