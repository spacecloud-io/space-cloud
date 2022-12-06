package graphql

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/graphql-go/graphql"

	"github.com/spacecloud-io/space-cloud/managers/source"
	graphqlProvider "github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func init() {
	caddy.RegisterModule(GraphqlSource{})
}

// GraphqlSource describes a graphql source
type GraphqlSource struct {
	v1alpha1.GraphqlSource

	// Internal stuff
	graphqlTypes map[string]graphql.Type

	queryRootType    graphql.Fields
	mutationRootType graphql.Fields
}

// CaddyModule returns the Caddy module information.
func (GraphqlSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(fmt.Sprintf("%s/v1alpha1", core.GroupName), "GraphqlSource")),
		New: func() caddy.Module { return new(GraphqlSource) },
	}
}

// Provision provisions the source
func (s *GraphqlSource) Provision(ctx caddy.Context) error {
	s.graphqlTypes = map[string]graphql.Type{
		graphql.Boolean.Name():  graphql.Boolean,
		graphql.String.Name():   graphql.String,
		graphql.Int.Name():      graphql.Int,
		graphql.Float.Name():    graphql.Float,
		graphql.ID.Name():       graphql.ID,
		graphql.DateTime.Name(): graphql.DateTime,
	}

	s.queryRootType = graphql.Fields{}
	s.mutationRootType = graphql.Fields{}

	// Fetch the graphql schema from the url
	if err := s.getSchemaFromSource(); err != nil {
		return err
	}

	return nil
}

// GetPriority returns the priority of the source. Higher
func (s *GraphqlSource) GetPriority() int {
	return 100
}

// GetTypes returns the root graphql types for this source
func (s *GraphqlSource) GetTypes() (queryTypes, mutationTypes graphql.Fields) {
	return s.queryRootType, s.mutationRootType
}

// GetAllTypes returns the all types in the source's type system
func (s *GraphqlSource) GetAllTypes() map[string]graphql.Type {
	return s.graphqlTypes
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
