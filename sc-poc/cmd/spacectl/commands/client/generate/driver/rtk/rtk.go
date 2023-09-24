package rtk

import (
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

type RTK struct {
	name string
}

func MakeRTKDriver(name string) driver.Driver {
	rtkDriver := &RTK{
		name: name,
	}
	return rtkDriver
}

func (r *RTK) GeneratePlugins([]v1alpha1.HTTPPlugin) (string, string, error) {
	return "", "", nil
}

var _ driver.Driver = (*RTK)(nil)
