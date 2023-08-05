package rtk

import "github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"

type RTK struct {
	name string
}

func MakeRTKDriver(name string) driver.Driver {
	rtkDriver := &RTK{
		name: name,
	}
	return rtkDriver
}

var _ driver.Driver = (*RTK)(nil)
