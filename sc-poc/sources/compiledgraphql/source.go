package compiledgraphql

import (
	"github.com/caddyserver/caddy/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/modules/rpc"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var compiledgraphqlsourcesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "compiledgraphqlsources"}

func init() {
	source.RegisterSource(Source{}, compiledgraphqlsourcesResource)
}

// Source describes the compiled graphql query source
type Source struct {
	v1alpha1.CompiledGraphqlSource

	// Internal stuff
	isReady       bool
	compiledQuery *graphql.CompiledQuery
}

// CaddyModule returns the Caddy module information.
func (Source) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(compiledgraphqlsourcesResource)),
		New: func() caddy.Module { return new(Source) },
	}
}

// Provision provisions the source
func (s *Source) Provision(ctx caddy.Context) error {
	return nil
}

// GetPriority returns the priority of the source. Higher
func (s *Source) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *Source) GetProviders() []string {
	providers := []string{"graphql"}
	if !s.Spec.InternalOnly {
		providers = append(providers, "rpc")
	}
	return providers
}

// Compiler resolves the graphql dependency
func (s *Source) GraphqlCompiler(fn graphql.CompilerFn) error {
	if s.isReady {
		return nil
	}

	// Mark the source as ready once compilation was successful
	defer func() {
		s.isReady = true
	}()
	return s.compile(fn)
}

// GetCompiledQuery returns the compiled query prepared
func (s *Source) GetCompiledQuery() *graphql.CompiledQuery {
	return s.compiledQuery
}

func (s *Source) GetRPCs() rpc.RPCs {
	requestSchema, responseSchema := s.getSchemas()
	r := &rpc.RPC{
		Name:          s.Name,
		OperationType: s.compiledQuery.OperationType,
		Extensions:    s.compiledQuery.Extensions,

		HTTPOptions: s.Spec.HTTP,
		Plugins:     s.Spec.Plugins,

		RequestSchema:  requestSchema,
		ResponseSchema: responseSchema,

		Call: s.call,
	}

	// Add the authentication function only if authentication is required
	if s.compiledQuery.IsAuthRequired {
		r.Authenticate = s.compiledQuery.AuthenticateRequest
	}

	return rpc.RPCs{r}
}

// Interface guards
var (
	_ caddy.Provisioner = (*Source)(nil)
	_ source.Source     = (*Source)(nil)
	_ rpc.Source        = (*Source)(nil)
	_ graphql.Compiler  = (*Source)(nil)
)
