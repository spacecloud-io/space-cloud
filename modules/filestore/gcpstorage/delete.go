package gcpstorage

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// DeleteFile deletes a file from GCPStorage
func (g *GCPStorage) DeleteFile(project, path string) error {
	// start of path must have single backslash
	// if end is a backslash, it will delete a folder
	if strings.HasSuffix(path, "/") {
		return g.DeleteDir(project, path)
	}
	return g.client.Bucket(project).Object(path).Delete(context.TODO())
}

// DeleteDir deletes a directory in GCPStorage
func (g *GCPStorage) DeleteDir(project, path string) error {
	// start and end of path must have single backslash
	bucket := g.client.Bucket(project)
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
