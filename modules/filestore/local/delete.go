package local

import (
	"strings"
	"os"
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
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}
