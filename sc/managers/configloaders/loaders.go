package configloaders

import (
	// importing file store module
	_ "github.com/spacecloud-io/space-cloud/managers/configloaders/file"

	// importing kube store module
	_ "github.com/spacecloud-io/space-cloud/managers/configloaders/kubernetes"
)
