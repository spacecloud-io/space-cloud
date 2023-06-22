package adapter

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
)

type ListOptions struct {
	Labels map[string]string `json:"labels"`
}

type Adapter interface {
	// Run starts the watcher.
	Run(context.Context) (chan common.ConfigType, error)

	// GetRawConfig returns the config in bytes.
	GetRawConfig() (common.ConfigType, error)

	// List returns all registered sources of a specific source type
	List(schema.GroupVersionResource, ListOptions) (*unstructured.UnstructuredList, error)

	// Get returns a registered source
	Get(schema.GroupVersionResource, string) (*unstructured.Unstructured, error)

	// Apply creates/updates a source
	Apply(schema.GroupVersionResource, *unstructured.Unstructured) error

	// Delete deletes a source
	Delete(schema.GroupVersionResource, string) error
}
