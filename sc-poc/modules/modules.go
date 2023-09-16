package modules

import (
	// Importing all applications
	_ "github.com/spacecloud-io/space-cloud/modules/auth"
	_ "github.com/spacecloud-io/space-cloud/modules/graphql"
	_ "github.com/spacecloud-io/space-cloud/modules/pubsub"
	_ "github.com/spacecloud-io/space-cloud/modules/rpc"
	_ "github.com/spacecloud-io/space-cloud/modules/tasks"
)
