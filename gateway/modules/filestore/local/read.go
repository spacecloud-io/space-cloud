package local

import (
	"bufio"
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ListDir lists the directory
func (l *Local) ListDir(ctx context.Context, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	ps := string(os.PathSeparator)
	path := strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(req.Path, ps)
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
func (l *Local) ReadFile(ctx context.Context, path string) (*model.File, error) {
	ps := string(os.PathSeparator)
	p := strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return &model.File{File: bufio.NewReader(f), Close: func() error { return f.Close() }}, nil
}
