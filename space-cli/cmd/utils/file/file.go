package file

import (
	"io/ioutil"
	"os"
)

type file interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
}

type def struct{}

// File is use for mocking
var File file

func init() {
	File = &def{}
}

// ReadFile is used to read a file d
func (d *def) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// WriteFile is used to write a file
func (d *def) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

// Stat is used to check if directory exits
func (d *def) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// IsNotExist is used to check if directory exits
func (d *def) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// MkdirAll is used to create a directory
func (d *def) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
