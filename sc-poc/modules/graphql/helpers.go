package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func (a *App) addToRootType(src *v1alpha1.GraphqlSource, rootType string, rootGraphqlType *graphql.Object) {
	a.logger.Debug("Attempting to load root type fields", zap.String("source", src.Name), zap.String("root_type", rootType))
	if t, p := a.graphqlTypes[fmt.Sprintf("%s_%s", src.Name, rootType)]; p {
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
			newFieldName := fmt.Sprintf("%s_%s", src.Name, fieldName)

			// Now lets merge the field with our root type
			rootGraphqlType.AddFieldConfig(newFieldName, &graphql.Field{
				Type:              field.Type,
				Args:              args,
				Description:       field.Description,
				DeprecationReason: field.DeprecationReason,
				Resolve:           a.resolveRemoteGraphqlQuery(src),
			})
			a.rootJoinObj[newFieldName] = ""
		}
		a.logger.Debug("Loaded root type fields", zap.String("source", src.Name), zap.String("root_type", rootType), zap.Int("# of field", len(q.Fields())), zap.Error(q.Error()))
	}
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

// TODO: Remove the following function
// func convertPrimitiveTypeToString(ogValue string, value interface{}) string {
// 	switch v := value.(type) {
// 	case (int):
// 		return strconv.Itoa(v)

// 	case (int64):
// 		return fmt.Sprintf("%v", v)

// 	case (float32):
// 		return fmt.Sprintf("%v", v)

// 	case (float64):
// 		return fmt.Sprintf("%v", v)

// 	case (string):
// 		return fmt.Sprintf("\"%s\"", v)

// 	case (bool):
// 		return strconv.FormatBool(v)

// 	default:
// 		return ogValue
// 	}
// }
