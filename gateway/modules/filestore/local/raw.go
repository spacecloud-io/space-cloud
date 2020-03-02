package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// DoesExists checks if the path exists
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

// GetState check if root path is valid
func (l *Local) GetState(ctx context.Context) error {
	if _, err := os.Stat(l.rootPath); os.IsNotExist(err) {
		return fmt.Errorf("root path does not exists")
	}
	return nil
}
