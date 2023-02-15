package k8s

import (
	"context"
	"fmt"

	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
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

// Get returns a registered source
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

// Apply creates/updates a source
func (k *K8s) Apply(gvr schema.GroupVersionResource, spec *unstructured.Unstructured) error {
	name := spec.GetName()
	oldObj, err := k.dc.Resource(gvr).Namespace(k.namespace).Get(context.TODO(), name, metav1.GetOptions{})
	// Resource not found. Create it
	if kubeErrors.IsNotFound(err) {
		if _, err := k.dc.Resource(gvr).Namespace(k.namespace).Create(context.TODO(), spec, metav1.CreateOptions{}); err != nil {
			return err
		}
		return nil
	}

	// Resource already exists. Update it
	spec.SetResourceVersion(oldObj.GetResourceVersion())
	if _, err := k.dc.Resource(gvr).Namespace(k.namespace).Update(context.TODO(), spec, metav1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

// Delete deletes a source
func (k *K8s) Delete(gvr schema.GroupVersionResource, name string) error {
	return k.dc.Resource(gvr).Namespace(k.namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
