package gcpstorage

import (
	"io"
	"strings"
	"context"

	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file in GCPStorage
func (g *GCPStorage) CreateFile(req *model.CreateFileRequest, file io.Reader) error {
	path := strings.Trim(req.Path, "/")
	name := strings.Trim(req.Name, "/")
	p := strings.Trim(path + "/" + name, "/")
	wc := g.client.Bucket(g.bucket).Object("/" + p).NewWriter(context.TODO())
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	return wc.Close()
}

// CreateDir creates a directory in GCPStorage
func (g *GCPStorage) CreateDir(req *model.CreateFileRequest) error {
	path := strings.Trim(req.Path, "/")
	name := strings.Trim(req.Name, "/")
	p := strings.Trim(path + "/" + name, "/")
	wc := g.client.Bucket(g.bucket).Object("/" + p + "/").NewWriter(context.TODO())
	_, err := wc.Write([]byte(""))
	if err != nil {
		return err
	}
	return wc.Close()
}
