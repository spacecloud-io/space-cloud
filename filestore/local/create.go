package local

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"

	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file in the path provided
func (l *Local) CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error {
	path := l.rootPath + "/" + project + req.Path
	if !isPathDir(path) {
		return errors.New("Local: Provided path is not a directory")
	}

	f, err := os.Create(path + "/" + req.Name)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, err = io.Copy(w, file)
	return err
}

// CreateDir creates a dirctory in the path provided
func (l *Local) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	path := l.rootPath + "/" + project + req.Path
	if !isPathDir(path) && !req.MakeAll {
		return errors.New("Local: Provided path is not a directory")
	}

	return os.MkdirAll(path+"/"+req.Name, os.ModePerm)
}
