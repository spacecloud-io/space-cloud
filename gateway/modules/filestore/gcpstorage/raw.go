package gcpstorage

import (
	"cloud.google.com/go/storage"
	"context"
)

// DoesExists checks if the path exists
func (g *GCPStorage) DoesExists(path string) error {
	if _, err := g.client.Bucket(g.bucket).Object(path).Attrs(context.TODO()); err != nil {
		return err
	}
	return nil
}

// GetState checks if sc is able to query gcp storage
func (g *GCPStorage) GetState() error {
	if _, err := g.client.Bucket(g.bucket).Object("/").Attrs(context.TODO()); err != nil && err != storage.ErrObjectNotExist {
		return err
	}
	return nil
}
