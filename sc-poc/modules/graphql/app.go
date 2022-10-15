package graphql

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("graphql", 100)
}

type App struct {
	// GraphqlSources contains the graphql sources to integrate with
	GraphqlSources []*v1alpha1.GraphqlSource `json:"graphqlSources"`

	// For internal use
	logger *zap.Logger

	// For graphql engine
	schema       graphql.Schema
	graphqlTypes map[string]graphql.Type

	rootQueryType *graphql.Object
	rootJoinObj   map[string]string
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "graphql",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the graphql module.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	return nil
}

// Start begins the graphql app operations
func (a *App) Start() error {
	// Create the root types
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})

	a.rootQueryType = queryType
	a.rootJoinObj = map[string]string{}

	// Create a new type map with preloaded types
	a.graphqlTypes = map[string]graphql.Type{
		graphql.Boolean.Name():  graphql.Boolean,
		graphql.String.Name():   graphql.String,
		graphql.Int.Name():      graphql.Int,
		graphql.Float.Name():    graphql.Float,
		graphql.ID.Name():       graphql.ID,
		graphql.DateTime.Name(): graphql.DateTime,
	}

	// Lets load the schemas for all sources
	for _, source := range a.GraphqlSources {
		queryRoot, mutationRoot, err := a.getSchemaFromUrl(source)
		if err != nil {
			a.logger.Error("Unable to get remote graphql schema", zap.String("source", source.Name), zap.Error(err))
			return err
		}

		// Extract the root types if provided
		a.addToRootType(source, queryRoot, queryType)
		a.addToRootType(source, mutationRoot, mutationType)
	}

	// Merge root types with schema if they are not empty
	schemaConfig := graphql.SchemaConfig{}
	if len(queryType.Fields()) > 0 {
		schemaConfig.Query = queryType
	}
	if len(mutationType.Fields()) > 0 {
		schemaConfig.Mutation = mutationType
	}

	// Add directives
	schemaConfig.Directives = []*graphql.Directive{
		{
			Name:      "export",
			Locations: []string{graphql.DirectiveLocationField},
			Args: []*graphql.Argument{
				{PrivateName: "as", Type: graphql.NewNonNull(graphql.String)},
			},
		}, {
			Name:      "auth",
			Locations: []string{graphql.DirectiveLocationSchema},
		}, {
			Name:      "injectClaim",
			Locations: []string{graphql.DirectiveLocationSchema, graphql.DirectiveLocationField},
			Args: []*graphql.Argument{
				{PrivateName: "key", Type: graphql.NewNonNull(graphql.String)},
				{PrivateName: "var", Type: graphql.String},
			},
		},
	}

	// Finally compile the graphql schema
	schmea, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		a.logger.Error("Unable to build schema object", zap.Error(err))
		return err
	}
	a.schema = schmea
	return nil
}

// Stop ends the graphql app operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
	_ apis.App          = (*App)(nil)
)
