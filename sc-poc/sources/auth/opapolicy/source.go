package opapolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/auth"
	authTypes "github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"

	"github.com/caddyserver/caddy/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var opapolicyResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "opapolicies"}

func init() {
	source.RegisterSource(OPAPolicySource{}, opapolicyResource)
}

// OPAPolicySource describes the OPAPolicy source
type OPAPolicySource struct {
	v1alpha1.OPAPolicy

	// For internal use
	logger                 *zap.Logger
	preparedQuery          rego.PreparedEvalQuery
	compiledGraphQLQueries map[string]*graphql.CompiledQuery
}

// CaddyModule returns the Caddy module information.
func (OPAPolicySource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(opapolicyResource)),
		New: func() caddy.Module { return new(OPAPolicySource) },
	}
}

// Provision provisions the source
func (s *OPAPolicySource) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)

	// Compile the rego policies
	if err := s.compileRegoPolicy(); err != nil {
		s.logger.Error("Unable to compile rego policies", zap.Error(err))
		return err
	}
	return nil
}

// GetPriority returns the priority of the source.
func (s *OPAPolicySource) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *OPAPolicySource) GetProviders() []string {
	return []string{"graphql", "auth"}
}

// Evaluate performs the evaluation of the policy and returns the result
func (s *OPAPolicySource) Evaluate(ctx context.Context, input interface{}) (bool, string, error) {
	rs, err := s.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, "", err
	}

	fmt.Println("OPA output---------------")
	d, _ := json.MarshalIndent(rs, "", " ")
	fmt.Println(string(d))
	fmt.Println("OPA output---------------")

	if len(rs) == 0 {
		return false, "Unable to evaluate opa policy", nil
	}

	if len(rs[0].Expressions) == 0 {
		return false, "Unable to evaluate opa policy", nil
	}

	// Extract output of opa policy
	result := authTypes.OPAOutput{
		Deny:   true,
		Allow:  false,
		Reason: "You are unauthorized to perform this request",
	}
	_ = mapstructure.Decode(rs[0].Bindings["result"], &result)

	// Check if request is allowed
	if result.Allow || !result.Deny {
		return true, "", nil
	}

	return false, result.Reason, nil
}

// SetCompiledQueries receives the compiledQueries from graphql app
func (s *OPAPolicySource) SetCompiledQueries(compiledQueries map[string]*graphql.CompiledQuery) {
	s.compiledGraphQLQueries = compiledQueries
}

func (s *OPAPolicySource) GetPluginDetails() v1alpha1.HTTPPlugin {
	return v1alpha1.HTTPPlugin{
		Name:   s.GetName(),
		Driver: "opa",
	}
}

// Interface guard
var (
	_ caddy.Provisioner             = (*OPAPolicySource)(nil)
	_ source.Source                 = (*OPAPolicySource)(nil)
	_ auth.PolicySource             = (*OPAPolicySource)(nil)
	_ graphql.CompiledQueryReceiver = (*OPAPolicySource)(nil)
	_ source.Plugin                 = (*OPAPolicySource)(nil)
)
