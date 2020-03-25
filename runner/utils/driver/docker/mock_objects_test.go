package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/stretchr/testify/mock"
	"github.com/txn2/txeh"
	"io"
)

type dockerClientMock struct {
	mock.Mock
}

func (a *dockerClientMock) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	args := a.Called(ctx, imageID)
	return args.Get(0).(types.ImageInspect), args.Get(1).([]byte), args.Error(2)
}
func (a *dockerClientMock) ImagePull(ctx context.Context, imageID string, options types.ImagePullOptions) (io.ReadCloser, error) {
	args := a.Called(ctx, imageID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (a *dockerClientMock) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{}, nil
}
func (a *dockerClientMock) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}
func (a *dockerClientMock) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return types.ContainerJSON{}, nil
}
func (a *dockerClientMock) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := a.Called(ctx, options)
	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}
	return value.([]types.Container), nil
}
func (a *dockerClientMock) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	args := a.Called(ctx, containerID, options)
	return args.Error(0)
}

type fileSystemMock struct {
	mock.Mock
}

func (f *fileSystemMock) RemoveTempSecretsFolder(projectID, serviceID, version string) error {
	return f.Called(projectID, serviceID, version).Error(0)
}

func (f *fileSystemMock) ReadSecretsFiles(ctx context.Context, projectID, secretName string) ([]byte, error) {
	args := f.Called(ctx, projectID, secretName)
	return args.Get(0).([]byte), args.Error(1)
}

func (f *fileSystemMock) CreateProjectDirectory(projectID string) error {
	return f.Called(projectID).Error(0)

}
func (f *fileSystemMock) RemoveProjectDirectory(projectID string) error {
	return f.Called(projectID).Error(0)
}

func (f *fileSystemMock) AddHostInHostFile(h *txeh.Hosts, IP, hostName string) {
	_ = f.Called(h, IP, hostName)
}

func (f *fileSystemMock) HostAddressLookUp(h *txeh.Hosts, hostName string) (bool, string, int) {
	args := f.Called(h, hostName)
	return args.Bool(0), args.String(1), args.Int(2)
}

func (f *fileSystemMock) SaveHostFile(h *txeh.Hosts) error {
	args := f.Called(h)
	return args.Error(0)
}

func (f *fileSystemMock) RemoveHostFromHostFile(h *txeh.Hosts, hostName string) {
	_ = f.Called(h, hostName)
}

func (f *fileSystemMock) NewHostFile() (*txeh.Hosts, error) {
	args := f.Called()
	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}
	logrus.Println("called")
	return value.(*txeh.Hosts), nil
}

type proxyMangerMock struct {
	mock.Mock
}

func (p *proxyMangerMock) SetServiceRoutes(projectID, serviceID string, r model.Routes) error {
	args := p.Called(projectID, serviceID, r)
	return args.Error(0)
}
func (p *proxyMangerMock) SetServiceRouteIfNotExists(projectID, serviceID, version string, ports []model.Port) error {
	args := p.Called(projectID, serviceID, version, ports)
	return args.Error(0)
}
func (p *proxyMangerMock) GetServiceRoutes(projectID string) (map[string]model.Routes, error) {
	args := p.Called(projectID)
	return args.Get(0).(map[string]model.Routes), args.Error(1)
}
func (p *proxyMangerMock) DeleteServiceRoutes(projectID, serviceID string) error {
	args := p.Called(projectID, serviceID)
	return args.Error(0)
}
