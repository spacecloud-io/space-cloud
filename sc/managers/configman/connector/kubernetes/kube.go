package kubernetes

import (
	"context"

	"github.com/spacecloud-io/space-cloud/model"
)

// Kube implements kube store
type Kube struct {
}

// ApplyResource applies resource in the store
func (f Kube) ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error {
	return nil
}

// GetResource gets resource from the store
func (f Kube) GetResource(ctx context.Context, meta *model.ResourceMeta) (*model.ResourceObject, error) {
	return nil, nil
}

// GetResources gets resources from the store
func (f Kube) GetResources(ctx context.Context, meta *model.ResourceMeta) (*model.ListResourceObjects, error) {
	return nil, nil
}

// DeleteResource delete resource from the store
func (f Kube) DeleteResource(ctx context.Context, meta *model.ResourceMeta) error {
	return nil
}

// DeleteResources delete resources from the store
func (f Kube) DeleteResources(ctx context.Context, meta *model.ResourceMeta) error {
	return nil
}
