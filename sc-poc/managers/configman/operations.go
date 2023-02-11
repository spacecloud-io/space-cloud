package configman

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List returns all the registered sources of a particular source type
func List(gvr schema.GroupVersionResource) ([]*unstructured.Unstructured, error) {
	return configLoader.adapter.List(gvr)
}
