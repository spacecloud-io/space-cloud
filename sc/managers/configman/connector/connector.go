package connector

import (
	"context"
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman/connector/file"
	"github.com/spacecloud-io/space-cloud/managers/configman/connector/kubernetes"

	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

// ConfigManConnector implemments store
type ConfigManConnector interface {
	ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error
	GetResource(ctx context.Context, meta *model.ResourceMeta) (*model.ResourceObject, error)
	GetResources(ctx context.Context, meta *model.ResourceMeta) (*model.ListResourceObjects, error)
	DeleteResource(ctx context.Context, meta *model.ResourceMeta) error
	DeleteResources(ctx context.Context, meta *model.ResourceMeta) error
}

// Connector connects stores
type Connector struct {
	Connector ConfigManConnector
}

// New returns a new store connector
func New(logger *zap.Logger, storeType, path string) (*Connector, error) {
	var connector Connector

	switch storeType {
	case "file":
		connector.Connector = file.File{Path: path}
	case "kube":
		connector.Connector = kubernetes.Kube{}
	default:
		return nil, fmt.Errorf("store-type: %s not supported", storeType)
	}

	return &connector, nil
}

// Destruct closes a store
func (c *Connector) Destruct() error {
	return nil
}
