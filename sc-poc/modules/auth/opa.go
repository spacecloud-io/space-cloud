package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"

	authTypes "github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/utils"
)

// EvaluateOPAPolicy evalues a stored opa query against the provided input
func (a *App) EvaluateOPAPolicy(ctx context.Context, name string, input interface{}) (bool, string, error) {
	query, p := a.regoPolicies[name]
	if !p {
		return false, "", fmt.Errorf("opa policy '%s' not found", name)
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
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

func (a *App) compileRegoPolicies() error {
	a.regoPolicies = make(map[string]rego.PreparedEvalQuery, len(a.OPAPolicies))
	for _, p := range a.OPAPolicies {

		// compiler, err := ast.CompileModulesWithOpt(map[string]string{p.Name: p.Spec.Rego}, ast.CompileOpts{})
		// if err != nil {
		// 	return fmt.Errorf("unable to compile rego policy '%s' - %s", p.Name, err.Error())
		// }
		query, err := a.compileRegoPolicy(p.Name, p.Spec.Rego)
		if err != nil {
			return err
		}
		a.regoPolicies[p.Name] = query
	}

	for _, q := range a.CompiledGraphqlSources {
		// Loop over all plugins
		for _, p := range q.Spec.Plugins {

			// We are only interested in plugins for the `auth_opa` driver
			if p.Driver != "auth_opa" {
				continue
			}

			// First lets get the params
			var params authTypes.PluginOPAParams
			_ = json.Unmarshal(p.Params.Raw, &params)

			// Compile the inline rego policy
			if params.Rego != "" {
				name := utils.Hash(params.Rego)
				query, err := a.compileRegoPolicy(name, params.Rego)
				if err != nil {
					return err
				}
				a.regoPolicies[name] = query
			}
		}
	}

	return nil
}

func (a *App) compileRegoPolicy(name, regoDoc string) (rego.PreparedEvalQuery, error) {
	r := rego.New(
		rego.Module(name, regoDoc),
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

				compiledQuery, p := a.graphqlApp.GetCompiledQuery(name)
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
	return r.PrepareForEval(context.TODO())
}
