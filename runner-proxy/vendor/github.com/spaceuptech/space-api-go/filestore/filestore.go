package filestore

import (
	"context"
	"io"
	"time"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Filestore contains the values for the filestore instance
type Filestore struct {
	config *config.Config
}

// New initializes the filestore module
func New(config *config.Config) *Filestore {
	return &Filestore{config}
}

// todo implement this
func (f *Filestore) CreateFolder(path, name string) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.CreateFolder(ctx, f.config.Project, path, name)
}

func (f *Filestore) DeleteFile(path string, meta interface{}) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.DeleteFile(ctx, meta, f.config.Project, path)
}

func (f *Filestore) ListFiles(listWhat, path string) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.List(ctx, f.config.Project, listWhat, path)
}

func (f *Filestore) UploadFile(path, name string, meta interface{}, reader io.Reader) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.UploadFile(ctx, f.config.Project, path, name, meta, reader)
}

func (f *Filestore) DownloadFile(path string, writer io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.DownloadFile(ctx, f.config.Project, path, writer)
}

func (f *Filestore) DoesFileOrFolderExists(path string) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.config.Transport.DoesExists(ctx, f.config.Project, path)
}
