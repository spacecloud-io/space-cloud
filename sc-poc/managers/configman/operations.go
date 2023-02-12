package configman

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all the registered sources of a particular source type
func List(gvr schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	return configLoader.adapter.List(gvr)
}

// Get returns a registered source
func Get(gvr schema.GroupVersionResource, name string) (*unstructured.Unstructured, error) {
	return configLoader.adapter.Get(gvr, name)
}
