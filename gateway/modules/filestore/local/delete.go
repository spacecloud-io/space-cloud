package local

import (
	"fmt"
	"os"
	"strings"
)

// DeleteDir deletes a directory if it exists
func (l *Local) DeleteDir(path string) error {
	ps := string(os.PathSeparator)
	path = strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	return os.RemoveAll(path)
}

// DeleteFile deletes a file if it exists
func (l *Local) DeleteFile(path string) error {
	ps := string(os.PathSeparator)
	path = strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	if isPathDir(path) {
		return fmt.Errorf("cannot delete the folder")
	}
	return os.Remove(path)
}
