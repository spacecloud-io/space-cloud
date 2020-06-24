package file

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type file interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Close(f *os.File) error
	Write(f *os.File, b []byte) (n int, err error)
	IsDir(f os.FileInfo) bool
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
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

// OpenFile is used to open a file
func (d *def) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// Close is used to close the file
func (d *def) Close(f *os.File) error {
	return f.Close()
}

// Write is use to write to a file
func (d *def) Write(f *os.File, b []byte) (n int, err error) {
	return f.Write(b)
}

// IsDir is use to check if path is directory
func (d *def) IsDir(f os.FileInfo) bool {
	return f.IsDir()
}

// Post makes http post request
func (d *def) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}
