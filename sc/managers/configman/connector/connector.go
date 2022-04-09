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
	SetLogger(logger *zap.Logger)
	Destruct() error
}

// New returns a new store connector
func New(logger *zap.Logger, storeType, path string) (ConfigManConnector, error) {
	switch storeType {
	case "file":
		return file.New(logger, path)
	case "kube":
		return kubernetes.New(logger)
	default:
		return nil, fmt.Errorf("store-type: %s not supported", storeType)
	}
}
