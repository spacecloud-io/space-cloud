package gcpstorage

import (
	"bufio"
	"io/ioutil"
	"os"
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	uuid "github.com/satori/go.uuid"
	
	"github.com/spaceuptech/space-cloud/model"
)

// ListDir lists a directory in GCPStorage
func (g *GCPStorage) ListDir(req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	// path should always start with a backslash
	path := strings.TrimRight(req.Path, "/") + "/"
	// path will always end in single backslash
	it := g.client.Bucket(g.bucket).Objects(context.TODO(), &storage.Query{
		Prefix: path,
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
	u2 := uuid.NewV4()

	tmpfile, err := ioutil.TempFile("", u2.String())
	if err != nil {
		return nil, err
	}

	rc, err := g.client.Bucket(g.bucket).Object(path).NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(tmpfile.Name(), data, 0644)
	if err != nil {
		return nil, err
	}
	
	tmpfile, err = os.Open(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	return &model.File{File: bufio.NewReader(tmpfile), Close: func() error {
		defer os.Remove(tmpfile.Name())
		return tmpfile.Close()
	}}, nil
}
