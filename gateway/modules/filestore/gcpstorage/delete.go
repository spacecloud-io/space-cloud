package gcpstorage

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// DeleteFile deletes a file from GCPStorage
func (g *GCPStorage) DeleteFile(path string) error {
	// start of path must have single backslash
	// if end is a backslash, it will delete a folder
	if strings.HasSuffix(path, "/") {
		return g.DeleteDir(path)
	}
	return g.client.Bucket(g.bucket).Object(path).Delete(context.TODO())
}

// DeleteDir deletes a directory in GCPStorage
func (g *GCPStorage) DeleteDir(path string) error {
	// start and end of path must have single backslash
	bucket := g.client.Bucket(g.bucket)
	it := bucket.Objects(context.TODO(), &storage.Query{
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
		err = bucket.Object(attrs.Name).Delete(context.TODO())
		if err != nil {
			return err
		}
	}
	return nil
}
