package configman

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all the registered sources of a particular source type
func List(gvr schema.GroupVersionResource, pkgName string) (*unstructured.UnstructuredList, error) {
	return configLoader.adapter.List(gvr, pkgName)
}

// Get returns a registered source
func Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	return configLoader.adapter.Get(gvr, name)
}

// Apply creates/updates a source
func Apply(gvr schema.GroupVersionResource, spec *unstructured.Unstructured) error {
	return configLoader.adapter.Apply(gvr, spec)
}

// Delete deletes a source
func Delete(gvr schema.GroupVersionResource, name string) error {
	return configLoader.adapter.Delete(gvr, name)
}
