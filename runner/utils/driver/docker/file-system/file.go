package file_system

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
	"io/ioutil"
	"os"
)

type fileSystem struct {
	secretPath string
}

func New(secretPath string) *fileSystem {
	return &fileSystem{secretPath: secretPath}
}

func (f *fileSystem) ReadSecretsFiles(ctx context.Context, projectID, secretName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s.json", f.secretPath, projectID, secretName))
	if err != nil {
		logrus.Errorf("unable to read file in file system %s - %v", secretName, err.Error())
		return nil, err
	}
	return data, err
}

func (f *fileSystem) RemoveTempSecretsFolder(projectID, serviceID, version string) error {
	return os.RemoveAll(fmt.Sprintf("%s/temp-secrets/%s/%s", f.secretPath, projectID, fmt.Sprintf("%s--%s", serviceID, version)))
}

func (f *fileSystem) CreateProjectDirectory(projectID string) error {
	projectPath := fmt.Sprintf("%s/%s", f.secretPath, projectID)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return os.MkdirAll(projectPath, 0777)
	}
	return nil
}

func (f *fileSystem) RemoveProjectDirectory(projectID string) error {
	return os.RemoveAll(fmt.Sprintf("%s/%s", f.secretPath, projectID))
}

func (f *fileSystem) SaveHostFile(h *txeh.Hosts) error {
	return h.Save()
}

func (f *fileSystem) RemoveHostFromHostFile(h *txeh.Hosts, hostName string) {
	h.RemoveHost(hostName)
}

func (f *fileSystem) AddHostInHostFile(h *txeh.Hosts, IP, hostName string) {
	h.AddHost(IP, hostName)
}

func (f *fileSystem) HostAddressLookUp(h *txeh.Hosts, hostName string) (bool, string, int) {
	return h.HostAddressLookup(hostName)
}

func (f *fileSystem) NewHostFile() (*txeh.Hosts, error) {
	return txeh.NewHostsDefault()
}
