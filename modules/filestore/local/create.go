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
	path := l.rootPath + project + req.Path

	// Create the dir recursively if it does not exists or overwrite if a file of same name already exists.
	if !isPathDir(path) {
		if !req.MakeAll {
			return errors.New("Local: Provided path is not a directory")
		}

		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path + string(os.PathSeparator) + req.Name)
	defer f.Close()
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = io.Copy(w, file)
	return err
}

// CreateDir creates a dirctory in the path provided
func (l *Local) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	path := l.rootPath + project + req.Path
	if !isPathDir(path) && !req.MakeAll {
		return errors.New("Local: Provided path is not a directory")
	}

	return os.MkdirAll(path+string(os.PathSeparator)+req.Name, os.ModePerm)
}
