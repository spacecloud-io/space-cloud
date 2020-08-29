package gcpstorage

import (
	"context"
	"io"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// CreateFile creates a file in GCPStorage
func (g *GCPStorage) CreateFile(ctx context.Context, req *model.CreateFileRequest, file io.Reader) error {
	req.Path = strings.TrimPrefix(req.Path, "/")
	path := req.Path + "/" + req.Name
	if len(req.Path) == 0 {
		path = req.Name
	}
	wc := g.client.Bucket(g.bucket).Object(path).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	return wc.Close()
}

// CreateDir creates a directory in GCPStorage
func (g *GCPStorage) CreateDir(ctx context.Context, req *model.CreateFileRequest) error {
	req.Path = strings.TrimPrefix(req.Path, "/")
	wc := g.client.Bucket(g.bucket).Object(utils.JoinTrailing(req.Path, req.Name, "/")).NewWriter(ctx)
	_, err := wc.Write([]byte(""))
	if err != nil {
		return err
	}
	return wc.Close()
}
