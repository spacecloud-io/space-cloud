package rootvalue

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
)

func extractQueryPrefixSuffix(sources []*types.Source, operation *ast.OperationDefinition, allowedVars map[string]struct{}, exportedVars map[string]*types.StoreValue) (prefix, suffix string) {

	// Add the operation and its value
	prefix = fmt.Sprintf("%s %s ", operation.Operation, operation.Name.Value)

	// Time to add the variables
	if len(operation.VariableDefinitions) > 0 {
		tempQuery := ""
		for _, v := range operation.VariableDefinitions {
			// Add the variable name
			varName := string(v.Variable.Loc.Source.Body[v.Variable.Loc.Start:v.Variable.Loc.End])

			// Check if the variable is allowed or not. We want to remove it if no one
			// is using this variable in the query
			if _, p := allowedVars[strings.TrimPrefix(varName, "$")]; !p {
				continue
			}

			tempQuery += varName
			// Add the type. Remember to trim the source name as well.
			tempQuery += ": "
			{
				temp := string(v.Type.GetLoc().Source.Body[v.Type.GetLoc().Start:v.Type.GetLoc().End])
				for _, src := range sources {
					if strings.HasPrefix(temp, src.Name) {
						temp = strings.TrimPrefix(temp, fmt.Sprintf("%s_", src.Name))
						break
					}
				}
				tempQuery += temp
			}

			// Add the default value
			if v.DefaultValue != nil {
				tempQuery += " = "
				tempQuery += string(v.DefaultValue.GetLoc().Source.Body[v.DefaultValue.GetLoc().Start:v.DefaultValue.GetLoc().End])
			}

			// Trailing comma
			tempQuery += ", "
		}

		// Now add the variable which have to be injected by dynamic variables created as a result of field exporting
		for k, v := range exportedVars {
			// Add the variable name
			tempQuery += "$" + k

			// Add the type. We don't need to trim the prefix here since it already is.
			tempQuery += ": " + v.TypeOf

			// Default value isn't needed since we'll always inject a variable for these kinds of variables

			// Trailing comma
			tempQuery += ", "
		}

		// Only add if there was atleast one qualified variable
		if len(tempQuery) > 0 {
			prefix += "(" + tempQuery + ") "
		}
	}

	for _, d := range operation.Directives {
		prefix += string(d.Loc.Source.Body[d.Loc.Start:d.Loc.End])
		prefix += " "
	}

	// Finish it off
	prefix += "{\n"
	suffix = "}\n"
	return
}
