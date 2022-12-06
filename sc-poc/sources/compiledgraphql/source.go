package compiledgraphql

import (
	"fmt"
	"sync"

	"github.com/caddyserver/caddy/v2"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/modules/rpc"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func init() {
	caddy.RegisterModule(Source{})
}

// Source describes the compiled graphql query source
type Source struct {
	v1alpha1.CompiledGraphqlSource

	// Internal stuff
	wg            *sync.WaitGroup
	compiledQuery *graphql.CompiledQuery
}

// CaddyModule returns the Caddy module information.
func (Source) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(fmt.Sprintf("%s/v1alpha1", core.GroupName), "CompiledGraphqlSource")),
		New: func() caddy.Module { return new(Source) },
	}
}

// Provision provisions the source
func (s *Source) Provision(ctx caddy.Context) error {
	s.wg = &sync.WaitGroup{}
	s.wg.Add(1)
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
	defer s.wg.Done()
	return s.compile(fn)
}

// GetCompiledQuery returns the compiled query prepared
func (s *Source) GetCompiledQuery() *graphql.CompiledQuery {
	return s.compiledQuery
}

func (s *Source) GetRPCs() rpc.RPCs {
	// Wait for the graphql configuration to be done
	s.wg.Wait()

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
