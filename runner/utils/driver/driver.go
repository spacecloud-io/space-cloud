package driver

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver/istio"
	v1 "k8s.io/api/core/v1"
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

		return istio.NewIstioDriver(auth, istioConfig)
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
}

// Driver is the interface of the modules which interact with the deployment targets
type Driver interface {
	CreateProject(project *model.Project) error
	ApplyService(service *model.Service) error
	AdjustScale(service *model.Service, activeReqs int32) error
	WaitForService(service *model.Service) error
	Type() model.DriverType

	// Secret methods!

	CreateSecret(secretObj *model.Secret, projectID string) error
	ListSecrets(secretObj *model.Secret, projectID string) (*v1.SecretList, error)
	DelSecrets(secretObj *model.Secret, projectID string) error
	SetKey(secretObj *model.SecretValue, projectID string, secretName string, secretKey string) error
	DelKey(projectID string, secretName string, secretKey string) error
}
