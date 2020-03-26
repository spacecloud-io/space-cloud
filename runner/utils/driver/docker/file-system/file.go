package filesystem

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
	"io/ioutil"
	"os"
)

// FileSystem stores information required for file operations
type FileSystem struct {
	secretPath string
}

// New creates a instance of file system
func New(secretPath string) *FileSystem {
	return &FileSystem{secretPath: secretPath}
}

// ReadSecretsFiles read files at path
func (f *FileSystem) ReadSecretsFiles(ctx context.Context, projectID, secretName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s.json", f.secretPath, projectID, secretName))
	if err != nil {
		logrus.Errorf("unable to read file in file system %s - %v", secretName, err.Error())
		return nil, err
	}
	return data, err
}

// RemoveTempSecretsFolder remove folder
func (f *FileSystem) RemoveTempSecretsFolder(projectID, serviceID, version string) error {
	return os.RemoveAll(fmt.Sprintf("%s/temp-secrets/%s/%s", f.secretPath, projectID, fmt.Sprintf("%s--%s", serviceID, version)))
}

// CreateProjectDirectory creates directory
func (f *FileSystem) CreateProjectDirectory(projectID string) error {
	projectPath := fmt.Sprintf("%s/%s", f.secretPath, projectID)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return os.MkdirAll(projectPath, 0777)
	}
	return nil
}

// RemoveProjectDirectory removes directory
func (f *FileSystem) RemoveProjectDirectory(projectID string) error {
	return os.RemoveAll(fmt.Sprintf("%s/%s", f.secretPath, projectID))
}

// SaveHostFile saves host file
func (f *FileSystem) SaveHostFile(h *txeh.Hosts) error {
	return h.Save()
}

// RemoveHostFromHostFile removes host from host file
func (f *FileSystem) RemoveHostFromHostFile(h *txeh.Hosts, hostName string) {
	h.RemoveHost(hostName)
}

// AddHostInHostFile adds hosts in host file
func (f *FileSystem) AddHostInHostFile(h *txeh.Hosts, IP, hostName string) {
	h.AddHost(IP, hostName)
}

// HostAddressLookUp check if provided host name is present in fiel
func (f *FileSystem) HostAddressLookUp(h *txeh.Hosts, hostName string) (bool, string, int) {
	return h.HostAddressLookup(hostName)
}

// NewHostFile creates new host file
func (f *FileSystem) NewHostFile() (*txeh.Hosts, error) {
	return txeh.NewHostsDefault()
}
