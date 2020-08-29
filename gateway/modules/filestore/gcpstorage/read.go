package gcpstorage

import (
	"bufio"
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ListDir lists a directory in GCPStorage
func (g *GCPStorage) ListDir(ctx context.Context, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	// path should not start with a backslash
	path := strings.Trim(req.Path, "/") + "/"
	if path == "/" {
		path = ""
	}

	it := g.client.Bucket(g.bucket).Objects(context.TODO(), &storage.Query{
		Prefix:    path,
		Delimiter: "/",
	})
	result := []*model.ListFilesResponse{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if attrs.Prefix != "" {
			prefix := strings.TrimPrefix(attrs.Prefix, req.Path)
			prefix = strings.TrimLeft(prefix, "/")
			t := &model.ListFilesResponse{Name: prefix, Type: "dir"}
			if req.Type == "all" || req.Type == t.Type {
				result = append(result, t)
			}
		} else {
			name := strings.TrimPrefix(attrs.Name, req.Path)
			name = strings.TrimLeft(name, "/")
			t := &model.ListFilesResponse{Name: name, Type: "file"}
			if req.Type == "all" || req.Type == t.Type {
				result = append(result, t)
			}
		}
	}
	return result, nil
}

// ReadFile reads a file from GCPStorage
func (g *GCPStorage) ReadFile(path string) (*model.File, error) {
	path = strings.TrimPrefix(path, "/")

	rc, err := g.client.Bucket(g.bucket).Object(path).NewReader(context.TODO())
	if err != nil {
		return nil, err
	}

	return &model.File{File: bufio.NewReader(rc), Close: func() error { return rc.Close() }}, nil
}
