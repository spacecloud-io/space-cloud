package graphql

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/managers/source"
)

func init() {
	caddy.RegisterModule(App{})
	apis.RegisterApp("graphql", 100)
}

type App struct {
	// For internal use
	logger *zap.Logger

	// For graphql engine
	schema       graphql.Schema
	graphqlTypes map[string]graphql.Type

	rootQueryType *graphql.Object
	rootJoinObj   map[string]string

	compiledQueries map[string]*CompiledQuery
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
	a.compiledQueries = make(map[string]*CompiledQuery)

	// Create a new type map with preloaded types
	a.graphqlTypes = map[string]graphql.Type{
		graphql.Boolean.Name():  graphql.Boolean,
		graphql.String.Name():   graphql.String,
		graphql.Int.Name():      graphql.Int,
		graphql.Float.Name():    graphql.Float,
		graphql.ID.Name():       graphql.ID,
		graphql.DateTime.Name(): graphql.DateTime,
	}

	sourceManT, err := ctx.App("source")
	if err != nil {
		a.logger.Error("Unable to load the source manager", zap.Error(err))
	}
	sourceMan := sourceManT.(*source.App)

	// Get all relevant sources
	sources := sourceMan.GetSources("graphql")

	// Iterate over all the sources to add them to the app
	for _, src := range sources {
		name := src.GetName()

		// First resolve the source's dependencies
		if err := source.ResolveDependencies(ctx, "graphql", src); err != nil {
			a.logger.Error("Unable to resolve source's dependency", zap.String("source", src.GetName()), zap.Error(err))
			return err
		}

		// Extract graphql types from the source
		graphqlSource, ok := src.(Source)
		if ok {
			queryFields, mutationFields := graphqlSource.GetTypes()
			a.addToRootType(name, queryType, queryFields, true)
			a.addToRootType(name, mutationType, mutationFields, false)

			// Extract all the types in this source's type system
			for k, v := range graphqlSource.GetAllTypes() {
				// Skip if we already have a types by this key
				if _, p := a.graphqlTypes[k]; p {
					continue
				}

				a.graphqlTypes[k] = v
			}
		}
	}

	// Merge root types with schema if they are not empty
	schemaConfig := graphql.SchemaConfig{}
	schemaConfig.Query = queryType
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
		}, {
			Name:      "tag",
			Locations: []string{graphql.DirectiveLocationSchema, graphql.DirectiveLocationField},
			Args: []*graphql.Argument{
				{PrivateName: "type", Type: graphql.NewNonNull(graphql.String)},
				{PrivateName: "key", Type: graphql.String},
			},
		},
	}

	// Finally compile the graphql schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		a.logger.Error("Unable to build schema object", zap.Error(err))
		return err
	}
	a.schema = schema

	// Time to process the dependant sources
	for _, src := range sources {
		name := src.GetName()

		dependantSource, ok := src.(Compiler)
		if ok {
			if err := dependantSource.GraphqlCompiler(a.Compile); err != nil {
				a.logger.Error("Unable to resolve dependencies for source", zap.String("name", name), zap.Error(err))
				return err
			}

			a.compiledQueries[name] = dependantSource.GetCompiledQuery()
		}

	}

	// Give compiledQueries to the receiver
	for _, src := range sources {
		receiver, ok := src.(CompiledQueryReceiver)
		if ok {
			receiver.SetCompiledQueries(a.compiledQueries)
		}
	}

	return nil
}

// Start begins the graphql app operations
func (a *App) Start() error {

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
