package graphql

import "github.com/graphql-go/graphql/language/ast"

func preprocessForAuth(graphqlDoc *ast.Document) (isAuthRequired bool, injectedClaims map[string]string) {
	// Set default values for result
	isAuthRequired = false
	injectedClaims = make(map[string]string)

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
		checkForClaimsDirective(s.(*ast.Field), injectedClaims)
	}

	return isAuthRequired || len(injectedClaims) > 0, injectedClaims
}

func checkForClaimsDirective(fieldAST *ast.Field, injectedClaims map[string]string) {
	// Loop over the directives
	for _, d := range fieldAST.Directives {
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
				checkForClaimsDirective(field, injectedClaims)
			}
		}
	}
}
