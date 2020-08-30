package gcpstorage

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
)

// DoesExists checks if the path exists
func (g *GCPStorage) DoesExists(ctx context.Context, path string) error {
	path = strings.TrimPrefix(path, "/")
	if _, err := g.client.Bucket(g.bucket).Object(path).Attrs(ctx); err != nil {
		return err
	}
	return nil
}

// GetState checks if sc is able to query gcp storage
func (g *GCPStorage) GetState(ctx context.Context) error {
	if _, err := g.client.Bucket(g.bucket).Object("/").Attrs(ctx); err != nil && err != storage.ErrObjectNotExist {
		return err
	}
	return nil
}
