package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
)

func (s *GraphqlSource) addToRootType(rootTypeName string, rootType graphql.Fields, graphqlTypes map[string]graphql.Type) {
	if t, p := graphqlTypes[fmt.Sprintf("%s_%s", s.Name, rootTypeName)]; p {
		q := t.(*graphql.Object)
		for fieldName, field := range q.Fields() {
			// Skip _join fields
			if fieldName == "_join" {
				continue
			}

			// Prepare args
			args := make(graphql.FieldConfigArgument, len(field.Args))
			for _, arg := range field.Args {
				args[arg.Name()] = &graphql.ArgumentConfig{
					Type:         arg.Type,
					Description:  arg.Description(),
					DefaultValue: arg.DefaultValue,
				}
			}

			// We need to change the field name to prevent conflicts with other fields
			newFieldName := fmt.Sprintf("%s_%s", s.Name, fieldName)

			// Now lets merge the field with our root type
			rootType[newFieldName] = &graphql.Field{
				Type:              field.Type,
				Args:              args,
				Description:       field.Description,
				DeprecationReason: field.DeprecationReason,
				Resolve:           s.resolveRemoteGraphqlQuery(),
			}

			// TODO: Add the type to the rootJoinObj in graphql provider
			// a.rootJoinObj[newFieldName] = ""
		}
	}
}

func extractQueryPrefixSuffix(operation *ast.OperationDefinition, allowedVars map[string]struct{}, exportedVars map[string]*types.StoreValue) (prefix, suffix string) {

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
				if strings.Contains(temp, "_") {
					arr := strings.Split(temp, "_")
					temp = strings.Join(arr[1:], "_")
				}

				// TODO: Remove this prefix removal logic
				// for _, src := range sources {
				// 	if strings.HasPrefix(temp, src.Name) {
				// 		temp = strings.TrimPrefix(temp, fmt.Sprintf("%s_", src.Name))
				// 		break
				// 	}
				// }
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
		// TODO: Skip all sc directives
		if d.Name.Value == "auth" {
			continue
		}

		prefix += string(d.Loc.Source.Body[d.Loc.Start:d.Loc.End])
		prefix += " "
	}

	// Finish it off
	prefix += "{\n"
	suffix = "}\n"
	return
}

func checkForVariablesInValue(value ast.Value, allowedVars map[string]struct{}) {
	switch v := value.(type) {
	case *ast.ObjectValue:
		for _, f := range v.Fields {
			checkForVariablesInValue(f.Value, allowedVars)
		}

	case *ast.ListValue:
		for _, f := range v.Values {
			checkForVariablesInValue(f, allowedVars)
		}

	case *ast.Variable:
		allowedVars[v.Name.Value] = struct{}{}
	}
}
