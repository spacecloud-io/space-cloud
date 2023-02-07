package opapolicy

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

// compileRegoPolicy compiles the OPA policy and stores
// the preparedQuery for future evaluation
func (s *OPAPolicySource) compileRegoPolicy() error {
	r := rego.New(
		rego.Module(s.Name, s.Spec.Rego),
		rego.Query("result := data[_]"),
		// rego.Compiler(compiler),
		rego.Function2(
			&rego.Function{
				Name:             "sc.executeCompiledQuery",
				Decl:             types.NewFunction(types.Args(types.S, types.A), types.A),
				Memoize:          true,
				Nondeterministic: true,
			},
			func(bctx rego.BuiltinContext, op1, op2 *ast.Term) (*ast.Term, error) {
				var name string
				var data map[string]interface{}

				// Parse the inputs first
				if err := ast.As(op1.Value, &name); err != nil {
					return nil, err
				}
				if err := ast.As(op2.Value, &data); err != nil {
					return nil, err
				}

				compiledQuery, p := s.compiledGraphQLQueries[name]
				if !p {
					return nil, fmt.Errorf("unable to get compiled query '%s'", name)
				}

				result := compiledQuery.Execute(bctx.Context, data)
				if result.HasErrors() {
					return nil, fmt.Errorf("unable to execute compiled graphql query - %s", result.Errors[0].Message)
				}

				v, err := ast.InterfaceToValue(result.Data)
				if err != nil {
					return nil, err
				}

				return ast.NewTerm(v), nil
			},
		),
	)
	preparedQuery, err := r.PrepareForEval(context.TODO())
	if err != nil {
		return err
	}

	s.preparedQuery = preparedQuery
	return nil
}

// TODO: compile inline rego policies in compiledGraphqlSources
