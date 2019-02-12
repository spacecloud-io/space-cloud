package local

import (
	"bufio"
	"context"
	"io/ioutil"
	"os"

	"github.com/spaceuptech/space-cloud/model"
)

// ListDir lists the directory
func (l *Local) ListDir(ctx context.Context, project string, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	path := l.rootPath + "/" + project + req.Path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := []*model.ListFilesResponse{}
	for _, f := range files {
		t := &model.ListFilesResponse{Name: f.Name(), Type: "file"}
		if f.IsDir() {
			t.Type = "dir"
		}

		if req.Type == "all" || req.Type == t.Type {
			result = append(result, t)
		}
	}

	return result, nil
}

// ReadFile reads a file from the path provided
func (l *Local) ReadFile(ctx context.Context, project, path string) (*model.File, error) {
	p := l.rootPath + "/" + project + path
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return &model.File{File: bufio.NewReader(f), Close: func() error { return f.Close() }}, nil
}
