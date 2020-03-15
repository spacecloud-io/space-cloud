package driver

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver/docker"
	"github.com/spaceuptech/space-cloud/runner/utils/driver/istio"
)

// New creates a new instance of the driver module
func New(auth *auth.Module, c *Config) (Driver, error) {

	switch c.DriverType {
	case model.TypeIstio:
		// Generate the config file
		var istioConfig *istio.Config
		if c.IsInCluster {
			istioConfig = istio.GenerateInClusterConfig()
		} else {
			istioConfig = istio.GenerateOutsideClusterConfig(c.ConfigFilePath)
		}
		istioConfig.SetProxyPort(c.ProxyPort)
		istioConfig.ArtifactAddr = c.ArtifactAddr

		return istio.NewIstioDriver(auth, istioConfig)

	case model.TypeDocker:
		return docker.NewDockerDriver(auth, c.ArtifactAddr)

	default:
		return nil, fmt.Errorf("invalid driver type (%s) provided", c.DriverType)
	}
}

// Config describes the configuration required by the driver module
type Config struct {
	DriverType     model.DriverType
	ConfigFilePath string
	IsInCluster    bool
	ProxyPort      uint32
	ArtifactAddr   string
}
type (
	// Driver is the interface of the modules which interact with the deployment targets
	Driver interface {
		CreateProject(ctx context.Context, project *model.Project) error
		DeleteProject(ctx context.Context, projectID string) error
		ApplyService(ctx context.Context, service *model.Service) error
		GetServices(ctx context.Context, projectID string) ([]*model.Service, error)
		DeleteService(ctx context.Context, projectID, serviceID, version string) error
		AdjustScale(service *model.Service, activeReqs int32) error
		WaitForService(service *model.Service) error
		Type() model.DriverType

		// Service routes

		ApplyServiceRoutes(ctx context.Context, projectID, serviceID string, routes model.Routes) error
		GetServiceRoutes(ctx context.Context, projectID string) (map[string]model.Routes, error)

		// Secret methods!

		CreateSecret(projectID string, secretObj *model.Secret) error
		ListSecrets(projectID string) ([]*model.Secret, error)
		DeleteSecret(projectID, secretName string) error
		SetKey(projectID, secretName, secretKey string, secretObj *model.SecretValue) error
		DeleteKey(projectID, secretName, secretKey string) error
		SetFileSecretRootPath(projectID string, secretName, rootPath string) error
	}
)
