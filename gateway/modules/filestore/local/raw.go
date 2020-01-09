package local

import (
	"errors"
	"os"
	"strings"
)

func (l *Local) DoesExists(path string) error {
	// check if file / folder exists
	ps := string(os.PathSeparator)
	path = strings.TrimRight(l.rootPath, ps) + ps + strings.TrimLeft(path, ps)
	if _, err := os.Stat(path); err != nil {
		// path does not exist
		return errors.New("provided file / dir path not found")
	}

	return nil
}
