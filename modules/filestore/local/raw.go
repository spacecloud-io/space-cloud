package local

import "os"

func (l *Local) DoesExists(path string) ( error) {
	// check if file / folder exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not exist
		return  nil
	}

	return  nil
}
