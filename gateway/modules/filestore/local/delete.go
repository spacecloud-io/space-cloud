package local

import (
	"context"
	"fmt"
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
		return fmt.Errorf("cannot delete the folder")
	}
	return os.Remove(path)
}
