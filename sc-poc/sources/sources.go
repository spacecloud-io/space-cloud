package sources

import (
	_ "github.com/spacecloud-io/space-cloud/sources/auth/jwthsasecret"
	_ "github.com/spacecloud-io/space-cloud/sources/auth/jwtrsasecret"
	_ "github.com/spacecloud-io/space-cloud/sources/auth/opapolicy"
	_ "github.com/spacecloud-io/space-cloud/sources/compiledgraphql"
	_ "github.com/spacecloud-io/space-cloud/sources/graphql"
	_ "github.com/spacecloud-io/space-cloud/sources/pubsub_channel"
	_ "github.com/spacecloud-io/space-cloud/sources/workspace"
)
