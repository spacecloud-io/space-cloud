package driver

import (
	"context"
	"fmt"
	"io"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver/istio"
)

// Config describes the configuration required by the driver module
type Config struct {
	DriverType     model.DriverType
	ConfigFilePath string
	IsInCluster    bool
	ProxyPort      uint32
	PrometheusAddr string
	ClusterName    string
}

// Interface is the interface of the modules which interact with the deployment targets
type Interface interface {
	CreateProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, projectID string) error
	ApplyService(ctx context.Context, service *model.Service) error
	GetServices(ctx context.Context, projectID string) ([]*model.Service, error)
	GetServiceStatus(ctx context.Context, projectID string) ([]*model.ServiceStatus, error)
	DeleteService(ctx context.Context, projectID, serviceID, version string) error
	ScaleUp(ctx context.Context, projectID, serviceID, version string) error
	WaitForService(ctx context.Context, service *model.Service) error
	Type() model.DriverType
	GetLogs(ctx context.Context, projectID string, info *model.LogRequest) (io.ReadCloser, error)

	// Service routes
	ApplyServiceRoutes(ctx context.Context, projectID, serviceID string, routes model.Routes) error
	GetServiceRoutes(ctx context.Context, projectID string) (map[string]model.Routes, error)

	// Service role
	ApplyServiceRole(ctx context.Context, role *model.Role) error
	GetServiceRole(ctx context.Context, projectID string) ([]*model.Role, error)
	DeleteServiceRole(ctx context.Context, projectID, serviceID, id string) error

	// Secret methods!
	CreateSecret(ctx context.Context, projectID string, secretObj *model.Secret) error
	ListSecrets(ctx context.Context, projectID string) ([]*model.Secret, error)
	DeleteSecret(ctx context.Context, projectID, secretName string) error
	SetKey(ctx context.Context, projectID, secretName, secretKey string, secretObj *model.SecretValue) error
	DeleteKey(ctx context.Context, projectID, secretName, secretKey string) error
	SetFileSecretRootPath(ctx context.Context, projectID string, secretName, rootPath string) error
}

// Module holds config of driver package
type Module struct {
	driver     Interface
	metricHook model.ServiceCallMetricHook
}

// New creates a new instance of the driver module
func New(auth *auth.Module, c *Config, hook model.ServiceCallMetricHook) (*Module, error) {
	d, err := initDriver(auth, c)
	if err != nil {
		return nil, err
	}
	return &Module{driver: d, metricHook: hook}, nil
}

func initDriver(auth *auth.Module, c *Config) (Interface, error) {
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
		istioConfig.PrometheusAddr = c.PrometheusAddr

		return istio.NewIstioDriver(auth, istioConfig)

	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("invalid driver type (%s) provided", c.DriverType), nil, nil)
	}
}
