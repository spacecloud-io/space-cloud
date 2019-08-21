package gcpstorage

import (
	"io"
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// CreateFile creates a file in GCPStorage
func (g *GCPStorage) CreateFile(req *model.CreateFileRequest, file io.Reader) error {
	wc := g.client.Bucket(g.bucket).Object(utils.JoinLeading(req.Path, req.Name, "/")).NewWriter(context.TODO())
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	return wc.Close()
}

// CreateDir creates a directory in GCPStorage
func (g *GCPStorage) CreateDir(req *model.CreateFileRequest) error {
	wc := g.client.Bucket(g.bucket).Object(utils.JoinLeadingTrailing(req.Path, req.Name, "/")).NewWriter(context.TODO())
	_, err := wc.Write([]byte(""))
	if err != nil {
		return err
	}
	return wc.Close()
}
