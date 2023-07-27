package golang

import "github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"

type Golang struct {
	pkgName string
}

func MakeGoDriver(pkgName string) driver.Driver {
	goDriver := &Golang{
		pkgName: pkgName,
	}
	return goDriver
}

var _ driver.Driver = (*Golang)(nil)
