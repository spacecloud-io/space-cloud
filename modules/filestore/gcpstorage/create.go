package gcpstorage

import (
	"io"
	"strings"
	"context"

	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file in GCPStorage
func (g *GCPStorage) CreateFile(project string, req *model.CreateFileRequest, file io.Reader) error {
	wc := g.client.Bucket(project).Object(req.Path + "/" + req.Name).NewWriter(context.TODO())
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	return wc.Close()
}

// CreateDir creates a directory in GCPStorage
func (g *GCPStorage) CreateDir(project string, req *model.CreateFileRequest) error {
	path := req.Path
	// back slash at the end is important, if not then file will be created of that name
	if !strings.HasSuffix(path, "/") {
		path = req.Path + "/"
	}

	wc := g.client.Bucket(project).Object(path).NewWriter(context.TODO())
	if _, err := wc.Write([]byte("")); err != nil {
		return err
	}
	return wc.Close()
}
