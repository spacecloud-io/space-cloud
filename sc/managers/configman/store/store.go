package connector

import (
	"context"
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/configman/store/file"
	"github.com/spacecloud-io/space-cloud/managers/configman/store/kubernetes"

	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

// Store implemments store
type Store interface {
	ApplyResource(ctx context.Context, resourceObj *model.ResourceObject) error
	GetResource(ctx context.Context, meta *model.ResourceMeta) (*model.ResourceObject, error)
	GetResources(ctx context.Context, meta *model.ResourceMeta) (*model.ListResourceObjects, error)
	DeleteResource(ctx context.Context, meta *model.ResourceMeta) error
	DeleteResources(ctx context.Context, meta *model.ResourceMeta) error
	SetLogger(logger *zap.Logger)
	Destruct() error
}

// New returns a new store connector
func New(logger *zap.Logger, storeType, path string) (Store, error) {
	switch storeType {
	case "file":
		return file.New(logger, path)
	case "kube":
		return kubernetes.New(logger)
	default:
		return nil, fmt.Errorf("store-type: %s not supported", storeType)
	}
}
