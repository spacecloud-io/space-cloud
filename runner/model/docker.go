package model

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/txn2/txeh"
	"io"
)

type DockerClient interface {
	ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
}

type FileSystem interface {
	ReadSecretsFiles(ctx context.Context, projectID, secretName string) ([]byte, error)
	RemoveTempSecretsFolder(projectID, serviceID, version string) error
	CreateProjectDirectory(projectID string) error
	RemoveProjectDirectory(projectID string) error

	SaveHostFile(h *txeh.Hosts) error
	RemoveHostFromHostFile(h *txeh.Hosts, hostName string)
	NewHostFile() (*txeh.Hosts, error)
	AddHostInHostFile(h *txeh.Hosts, IP, hostName string)
	HostAddressLookUp(h *txeh.Hosts, hostName string) (bool, string, int)
}

type ProxyManager interface {
	SetServiceRoutes(projectID, serviceID string, r Routes) error
	SetServiceRouteIfNotExists(projectID, serviceID, version string, ports []Port) error
	GetServiceRoutes(projectID string) (map[string]Routes, error)
	DeleteServiceRoutes(projectID, serviceID string) error
}
