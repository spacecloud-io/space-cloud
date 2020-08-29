package local

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// CreateFile creates a file in the path provided
func (l *Local) CreateFile(ctx context.Context, req *model.CreateFileRequest, file io.Reader) error {
	ps := string(os.PathSeparator)
	path := strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(req.Path, ps)

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
	defer utils.CloseTheCloser(f)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer func() { _ = w.Flush() }()

	_, err = io.Copy(w, file)
	return err
}

// CreateDir creates a directory in the path provided
func (l *Local) CreateDir(ctx context.Context, req *model.CreateFileRequest) error {
	ps := string(os.PathSeparator)
	path := strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(req.Path, ps)
	if !isPathDir(path) && !req.MakeAll {
		return errors.New("Local: Provided path is not a directory")
	}

	return os.MkdirAll(path+string(os.PathSeparator)+req.Name, os.ModePerm)
}
