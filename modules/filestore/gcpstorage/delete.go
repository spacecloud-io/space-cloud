package gcpstorage

import (
	"context"
)

// DeleteFile deletes a file from GCPStorage
func (g *GCPStorage) DeleteFile(project, path string) error {
	return g.client.Bucket(project).Object(path).Delete(context.TODO())
}

// DeleteDir deletes a directory in GCPStorage
func (g *GCPStorage) DeleteDir(project, path string) error {
	return g.client.Bucket(project).Object(path).Delete(context.TODO())
}
