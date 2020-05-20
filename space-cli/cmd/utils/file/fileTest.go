package file

import (
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
	return nil, c.Error(1)
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
