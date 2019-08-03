package gcpstorage

import (
	"context"

	"cloud.google.com/go/storage"

	"github.com/spaceuptech/space-cloud/utils"
)

// GCPStorage holds the GCPStorage client
type GCPStorage struct {
	client *storage.Client
	bucket string
}

// Init initializes a GCPStorage client
func Init(bucket string) (*GCPStorage, error) {
	ctx := context.TODO()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCPStorage{client, bucket}, nil
}

// GetStoreType returns the file store type
func (g *GCPStorage) GetStoreType() utils.FileStoreType {
	return utils.GCPStorage
}

// Gracefully close the GCPStorage module
func (g *GCPStorage) Close() error {
	return g.Close()
}
