package gcpstorage

import (
	"fmt"
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"context"
	// "strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	uuid "github.com/satori/go.uuid"
	
	"github.com/spaceuptech/space-cloud/model"
)

// ListDir lists a directory in GCPStorage
func (g *GCPStorage) ListDir(project string, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	it := g.client.Bucket(project).Objects(context.TODO(), &storage.Query{
		Prefix:    req.Path, //backslash at the end is important
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
		fmt.Println(attrs)
		if attrs.Prefix != "" {
			t := &model.ListFilesResponse{Name: filepath.Base(attrs.Name), Type: "file"}
			if req.Type == "all" || req.Type == t.Type {
				result = append(result, t)
			}
		} else {
			t := &model.ListFilesResponse{Name: filepath.Base(attrs.Prefix), Type: "dir"}
			if req.Type == "all" || req.Type == t.Type {
				result = append(result, t)
			}
		}
	}
	return result, nil
}

// ReadFile reads a file from GCPStorage
func (g *GCPStorage) ReadFile(project, path string) (*model.File, error) {
	u2 := uuid.NewV4()

	tmpfile, err := ioutil.TempFile("", u2.String())
	if err != nil {
		return nil, err
	}

	rc, err := g.client.Bucket(project).Object(project).NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	if _, err := tmpfile.Write(data); err != nil {
		return nil, err
	}

	return &model.File{File: bufio.NewReader(tmpfile), Close: func() error {
		defer os.Remove(tmpfile.Name())
		return tmpfile.Close()
	}}, nil
}
