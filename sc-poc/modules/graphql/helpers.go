package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"github.com/spacecloud-io/space-cloud/utils"
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

func (a *App) convertGraphqlOutputToJSONSchema(operationAST *ast.OperationDefinition) (*jsonschema.Schema, map[string]interface{}, error) {
	required := []string{}
	properties := orderedmap.New()

	extensions := make(map[string]interface{})

	// Get operation fields
	var fields graphql.FieldDefinitionMap
	if operationAST.Operation == "query" {
		fields = a.schema.QueryType().Fields()
	} else {
		fields = a.schema.MutationType().Fields()
	}

	for _, field := range operationAST.SelectionSet.Selections {
		fieldAST := field.(*ast.Field)
		fieldName := fieldAST.Name.Value
		propertyName := getFieldNameOrAlias(fieldAST)

		fieldDef, p := fields[fieldName]
		if !p {
			return nil, nil, fmt.Errorf("invalid graphql type '%s' provided for field", fieldName)
		}

		schema, isRequired := a.convertFieldTypeToJSONSchema(fmt.Sprintf("/%s", propertyName), fieldDef.Type, fieldAST, extensions)
		if isRequired {
			required = append(required, propertyName)
		}
		properties.Set(propertyName, schema)
	}

	return &jsonschema.Schema{Type: "object", Properties: properties, Required: required}, extensions, nil
}

func (a *App) convertFieldTypeToJSONSchema(path string, t graphql.Type, fieldAST *ast.Field, extensions map[string]interface{}) (*jsonschema.Schema, bool) {
	switch v := t.(type) {
	case *graphql.Scalar, *graphql.Enum:
		// Add exposed types in the extensions
		insertTagInExtensions(path, fieldAST, extensions)

		return a.convertGraphqlTypeToJSONSchema("", v, map[string]struct{}{})
	case *graphql.List:
		js, _ := a.convertFieldTypeToJSONSchema(path, v.OfType, fieldAST, extensions)
		return &jsonschema.Schema{Type: "array", Items: js}, false
	case *graphql.NonNull:
		js, _ := a.convertFieldTypeToJSONSchema(path, v.OfType, fieldAST, extensions)
		return js, true
	case *graphql.Object:
		// Add exposed types in the extensions
		insertTagInExtensions(path, fieldAST, extensions)

		required := []string{}
		properties := orderedmap.New()

		fields := v.Fields()

		for _, sel := range fieldAST.SelectionSet.Selections {
			subFieldAST := sel.(*ast.Field)
			fieldDef := fields[subFieldAST.Name.Value]

			// Name is the property propertyName the client will receive
			propertyName := getFieldNameOrAlias(subFieldAST)

			schema, isRequired := a.convertFieldTypeToJSONSchema(fmt.Sprintf("%s/%s", path, propertyName), fieldDef.Type, subFieldAST, extensions)
			if isRequired {
				required = append(required, propertyName)
			}

			properties.Set(propertyName, schema)
		}

		return &jsonschema.Schema{Type: "object", Properties: properties, Required: required}, false
	}

	return nil, false
}

func (a *App) convertVariablesToJSONSchema(operationAST *ast.OperationDefinition, ignoredFields map[string]struct{}) (*jsonschema.Schema, error) {
	required := []string{}
	properties := orderedmap.New()
	// Iterate over all the variables
	for _, v := range operationAST.VariableDefinitions {
		varName := v.Variable.Name.Value

		// Skip the field if it's to be ignored
		if _, p := ignoredFields[varName]; p {
			continue
		}

		varType, isList, isRequired := getGraphqlTypeNameFromAST(v.Type)
		graphqlType, p := a.graphqlTypes[varType]
		if !p {
			return nil, fmt.Errorf("invalid graphql type '%s' provided for variable '%s'", varType, varName)
		}

		js, _ := a.convertGraphqlTypeToJSONSchema(varType, graphqlType, map[string]struct{}{})

		// Check if field is required
		if isRequired {
			required = append(required, varName)
		}

		// Check if field is an array
		if isList {
			js = &jsonschema.Schema{Type: "array", Items: js}
		}

		properties.Set(varName, js)
	}

	return &jsonschema.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}, nil
}

func (a *App) convertGraphqlTypeToJSONSchema(propertyType string, t graphql.Type, capturedFields map[string]struct{}) (*jsonschema.Schema, bool) {
	capturedFields[propertyType] = struct{}{}
	switch v := t.(type) {
	case *graphql.Scalar:
		switch v.Name() {
		case graphql.ID.Name(), graphql.DateTime.Name(), graphql.String.Name():
			return &jsonschema.Schema{Type: "string"}, false
		case graphql.Boolean.Name():
			return &jsonschema.Schema{Type: "boolean"}, false
		case graphql.Int.Name():
			return &jsonschema.Schema{Type: "integer"}, false
		case graphql.Float.Name():
			return &jsonschema.Schema{Type: "number"}, false
		default:
			return &jsonschema.Schema{Type: "string"}, false
		}

	case *graphql.InputObject:

		required := []string{}
		properties := orderedmap.New()
		for _, f := range v.Fields() {
			fieldType := f.Type.Name()
			if _, p := capturedFields[fieldType]; p {
				continue
			}

			js, isRequired := a.convertGraphqlTypeToJSONSchema(fieldType, f.Type, utils.ShallowCopy(capturedFields))
			if isRequired {
				required = append(required, f.Name())
			}

			properties.Set(f.Name(), js)
		}

		return &jsonschema.Schema{Type: "object", Properties: properties, Required: required}, false

		// _, isValidGraphqlType := a.graphqlTypes[propertyType]
		// _, isPresentInDef := defs[propertyType]

		// schema := new(jsonschema.Schema)
		// if !isValidGraphqlType || !isPresentInDef {
		// 	// Add type in definitions if its a valid graphql type
		// 	if isValidGraphqlType {
		// 		defs[propertyType] = schema
		// 	}

		// 	required := []string{}
		// 	properties := orderedmap.New()
		// 	for _, f := range v.Fields() {
		// 		js, isRequired := a.convertGraphqlTypeToJSONSchema(f.Type.Name(), f.Type, defs)
		// 		if isRequired {
		// 			required = append(required, f.Name())
		// 		}

		// 		properties.Set(f.Name(), js)
		// 	}

		// 	schema.Type = "object"
		// 	schema.Properties = properties
		// 	schema.Required = required

		// 	if !isValidGraphqlType {
		// 		return schema, false
		// 	}
		// }

		// return &jsonschema.Schema{Ref: fmt.Sprintf("#/$defs/%s", propertyType)}, false

	// case *graphql.Object:
	// 	required := []string{}
	// 	properties := orderedmap.New()
	// 	for _, f := range v.Fields() {
	// 		if _, p := capturedType[f.Name]; p {
	// 			continue
	// 		}
	// 		capturedType[f.Name] = struct{}{}

	// 		js, isRequired := convertGraphqlTypeToJSONSchema(f.Type, capturedType)
	// 		if isRequired {
	// 			required = append(required, f.Name)
	// 		}

	// 		properties.Set(f.Name, js)
	// 	}
	// 	return &jsonschema.Schema{Type: "object", Properties: properties, Required: required}, false

	// case *graphql.Argument:
	// 	return convertGraphqlTypeToJSONSchema(v. v.Type, capturedType)

	case *graphql.Enum:
		values := make([]interface{}, len(v.Values()))
		for i, item := range v.Values() {
			values[i] = item.Value
		}
		return &jsonschema.Schema{Type: "string", Enum: values}, false

	case *graphql.List:
		js, _ := a.convertGraphqlTypeToJSONSchema(v.OfType.Name(), v.OfType, capturedFields)
		return &jsonschema.Schema{Type: "array", Items: js}, false

	case *graphql.NonNull:
		js, _ := a.convertGraphqlTypeToJSONSchema(v.OfType.Name(), v.OfType, capturedFields)
		return js, true

	default:
		return nil, false
	}
}

func getFieldNameOrAlias(fieldAST *ast.Field) string {
	if fieldAST.Alias != nil {
		return fieldAST.Alias.Value
	}

	return fieldAST.Name.Value
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

func insertTagInExtensions(path string, fieldAST *ast.Field, extensions map[string]interface{}) {
	for _, d := range fieldAST.Directives {
		if d.Name.Value != "tag" {
			continue
		}

		var t, k string
		for _, a := range d.Arguments {
			switch a.Name.Value {
			case "type":
				t = a.Value.(*ast.StringValue).Value
			case "key":
				k = a.Value.(*ast.StringValue).Value
			}
		}

		// Use the path as the key if it wasn't provided
		if k == "" {
			k = path
		}

		var arr []interface{}
		if t, p := extensions["x-tags"]; p {
			arr = t.([]interface{})
		}
		extensions["x-tags"] = append(arr, map[string]interface{}{"type": t, "key": k})
		break
	}
}
