package graphql

import "github.com/graphql-go/graphql/language/ast"

func preprocessForAuth(graphqlDoc *ast.Document) (isAuthRequired bool, injectedClaims map[string]string, exportedVars map[string]struct{}) {
	// Set default values for result
	isAuthRequired = false
	injectedClaims = make(map[string]string)
	exportedVars = make(map[string]struct{})

	// Get the operation ast
	operationAST := graphqlDoc.Definitions[0].(*ast.OperationDefinition)

	// Check authentication directive is present in the schema
	for _, d := range operationAST.Directives {
		if d.Name.Value == "auth" {
			isAuthRequired = true
			break
		}
	}

	// Parse each field to check if any claims are to be injected
	for _, s := range operationAST.SelectionSet.Selections {
		checkForInjectedAndExportedVars(s.(*ast.Field), injectedClaims, exportedVars)
	}

	return isAuthRequired || len(injectedClaims) > 0, injectedClaims, exportedVars
}

func checkForInjectedAndExportedVars(fieldAST *ast.Field, injectedClaims map[string]string, exportedVars map[string]struct{}) {
	// Loop over the directives
	for _, d := range fieldAST.Directives {
		if d.Name.Value == "export" {
			as := d.Arguments[0].Value.(*ast.StringValue).Value
			exportedVars[as] = struct{}{}
		}

		if d.Name.Value == "injectClaim" {
			// Find the key and variable
			var key, variable string
			for _, arg := range d.Arguments {
				if arg.Name.Value == "key" {
					key = arg.Value.(*ast.StringValue).Value
				}
				if arg.Name.Value == "var" {
					variable = arg.Value.(*ast.StringValue).Value
				}
			}

			if variable == "" {
				variable = key
			}
			injectedClaims[key] = variable
		}
	}

	// Do the same for each of the selection set
	if fieldAST.SelectionSet != nil {
		for _, s := range fieldAST.SelectionSet.Selections {
			if field, ok := s.(*ast.Field); ok {
				checkForInjectedAndExportedVars(field, injectedClaims, exportedVars)
			}
		}
	}
}
