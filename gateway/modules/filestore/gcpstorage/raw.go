package gcpstorage

import (
	"context"
)

func (g *GCPStorage) DoesExists(path string)  error{
	if _, err := g.client.Bucket(g.bucket).Object(path).Attrs(context.TODO()); err != nil {
		return  err
	}
	return  nil
}
