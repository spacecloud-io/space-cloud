package local

import (
	"context"
	"os"
	"strings"
)

// DeleteDir deletes a directory if it exists
func (l *Local) DeleteDir(ctx context.Context, path string) error {
	ps := string(os.PathSeparator)
	path = strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	return os.RemoveAll(path)
}

// DeleteFile deletes a file if it exists
func (l *Local) DeleteFile(ctx context.Context, path string) error {
	ps := string(os.PathSeparator)
	path = strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	if isPathDir(path) {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}
