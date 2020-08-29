package gcpstorage

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// DeleteFile deletes a file from GCPStorage
func (g *GCPStorage) DeleteFile(ctx context.Context, path string) error {
	// trim / at the start and at the end
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	return g.client.Bucket(g.bucket).Object(path).Delete(ctx)
}

// DeleteDir deletes a directory in GCPStorage
func (g *GCPStorage) DeleteDir(ctx context.Context, path string) error {
	path = strings.TrimPrefix(path, "/")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	bucket := g.client.Bucket(g.bucket)
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: path,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		err = bucket.Object(attrs.Name).Delete(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
