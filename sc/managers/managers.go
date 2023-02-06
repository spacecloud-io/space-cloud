package manager

import (
	// Importing config loader module
	_ "github.com/spacecloud-io/space-cloud/managers/configloaders"

	// Importing config handler module
	_ "github.com/spacecloud-io/space-cloud/managers/configman"

	// Importing api handler module
	_ "github.com/spacecloud-io/space-cloud/managers/apis"
)
