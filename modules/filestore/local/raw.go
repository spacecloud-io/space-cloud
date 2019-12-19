package local

import (
	"errors"
	"os"
)

func (l *Local) DoesExists(path string) error {
	// check if file / folder exists
	if _, err := os.Stat(path); err != nil {
		// path does not exist
		return errors.New("provided file / dir path not found")
	}

	return nil
}
