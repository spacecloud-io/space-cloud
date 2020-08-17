package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/stretchr/testify/mock"
)

// Mocket used during test
type Mocket struct {
	mock.Mock
}

// ReadFile is used to read a file during test
func (m *Mocket) ReadFile(filename string) ([]byte, error) {
	c := m.Called()
	return c.Get(0).([]byte), c.Error(1)
}

// WriteFile is used to write a file during test
func (m *Mocket) WriteFile(filename string, data []byte, perm os.FileMode) error {
	c := m.Called()
	return c.Error(0)
}

// Stat is used to check if directory exits during test
func (m *Mocket) Stat(name string) (os.FileInfo, error) {
	c := m.Called()
	var x os.FileInfo
	return x, c.Error(1)
}

// IsNotExist is used to check if directory exits during test
func (m *Mocket) IsNotExist(err error) bool {
	c := m.Called()
	return c.Bool(0)
}

// MkdirAll is used to create a directory during test
func (m *Mocket) MkdirAll(path string, perm os.FileMode) error {
	c := m.Called()
	return c.Error(0)
}

// OpenFile is used to create a directory during test
func (m *Mocket) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	c := m.Called()
	return &os.File{}, c.Error(1)
}

// Close is used to close the file during test
func (m *Mocket) Close(f *os.File) error {
	c := m.Called()
	return c.Error(0)
}

// Write is use to write to a file during test
func (m *Mocket) Write(f *os.File, b []byte) (n int, err error) {
	c := m.Called()
	return c.Int(0), c.Error(1)
}

// IsDir is use to check if path is directory during test
func (m *Mocket) IsDir(f os.FileInfo) bool {
	c := m.Called()
	return c.Bool(0)
}

// Post makes http post request during testing
func (m *Mocket) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	c := m.Called()
	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	return &http.Response{StatusCode: c.Int(0), Body: r}, c.Error(1)
}
