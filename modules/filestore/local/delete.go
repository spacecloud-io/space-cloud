package local

import (
	"context"
	"errors"
	"os"
)

// DeleteDir deletes a directory if it exists
func (l *Local) DeleteDir(ctx context.Context, project, path string) error {
	p := l.rootPath + "/" + project + path

	return os.RemoveAll(p)
}

// DeleteFile deletes a directory if it exists
func (l *Local) DeleteFile(ctx context.Context, project, path string) error {
	p := l.rootPath + "/" + project + path
	if isPathDir(p) {
		return errors.New("Local: Provided path is not a directory")
	}

	return os.Remove(p)
}
