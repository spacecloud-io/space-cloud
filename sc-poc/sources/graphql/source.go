package graphql

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/graphql-go/graphql"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/source"
	graphqlProvider "github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var graphqlsourcesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "graphqlsources"}

func init() {
	source.RegisterSource(GraphqlSource{}, graphqlsourcesResource)
}

// GraphqlSource describes a graphql source
type GraphqlSource struct {
	v1alpha1.GraphqlSource

	// Internal stuff
	rawSchema *introspectionResponse
}

// CaddyModule returns the Caddy module information.
func (GraphqlSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(graphqlsourcesResource)),
		New: func() caddy.Module { return new(GraphqlSource) },
	}
}

// Provision provisions the source
func (s *GraphqlSource) Provision(ctx caddy.Context) error {
	return s.getRawGraphqlSchema()
}

// GetPriority returns the priority of the source. Higher
func (s *GraphqlSource) GetPriority() int {
	return 100
}

// GetTypes returns the root graphql types for this source
func (s *GraphqlSource) GetGraphQLTypes() *graphqlProvider.Types {
	graphqlTypes := map[string]graphql.Type{
		graphql.Boolean.Name():  graphql.Boolean,
		graphql.String.Name():   graphql.String,
		graphql.Int.Name():      graphql.Int,
		graphql.Float.Name():    graphql.Float,
		graphql.ID.Name():       graphql.ID,
		graphql.DateTime.Name(): graphql.DateTime,
	}

	queryRootType := graphql.Fields{}
	mutationRootType := graphql.Fields{}

	// Fetch the graphql schema from the url
	s.prepareGraphqlTypes(queryRootType, mutationRootType, graphqlTypes)

	return &graphqlProvider.Types{
		QueryTypes:    queryRootType,
		MutationTypes: mutationRootType,
		AllTypes:      graphqlTypes,
	}
}

// GetProviders returns the providers this source is applicable for
func (s *GraphqlSource) GetProviders() []string {
	return []string{"graphql"}
}

// Interface guards
var (
	_ caddy.Provisioner      = (*GraphqlSource)(nil)
	_ source.Source          = (*GraphqlSource)(nil)
	_ graphqlProvider.Source = (*GraphqlSource)(nil)
)
