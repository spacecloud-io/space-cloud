package local

import (
	"os"

	"github.com/spaceuptech/space-cloud/utils"
)

// Local is the file store driver for the local filesystem
type Local struct {
	rootPath string
}

// Init initialises the local filestore driver
func Init(path string) (*Local, error) {
	return &Local{path}, nil
}

// GetStoreType returns the file store type
func (l *Local) GetStoreType() utils.FileStoreType {
	return utils.Local
}

// Close gracefully closed the local filestore module
func (l *Local) Close() error {
	return nil
}

func isPathDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}
