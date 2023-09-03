package typescript

import "github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"

type Typescript struct {
}

func MakeTSDriver() driver.Driver {
	tsDriver := &Typescript{}
	return tsDriver
}

var _ driver.Driver = (*Typescript)(nil)
