package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"

	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) addToRootType(srcName string, rootGraphqlType *graphql.Object, graphqlType graphql.Fields, injectJoin bool) {
	for fieldName, field := range graphqlType {
		// Skip _join fields
		if fieldName == "_join" {
			continue
		}

		// Add the missing resolvers
		a.addResolvers(srcName, field.Type, injectJoin)

		// Add the field to the root graphql type
		rootGraphqlType.AddFieldConfig(fieldName, field)

		// Add the field to the root join object
		if injectJoin {
			a.rootJoinObj[fieldName] = ""
		}
	}
}

func (a *App) addResolvers(srcName string, graphqlType graphql.Output, injectJoin bool) {
	graphqlObject, ok := graphqlType.(*graphql.Object)
	if !ok {
		return
	}

	// Get all the fields of this type
	fields := graphqlObject.Fields()

	// Skip if join field is already present
	if _, p := fields["_join"]; p {
		return
	}

	for _, field := range fields {
		nestedGraphqlObject, ok := field.Type.(*graphql.Object)
		if ok {
			a.addResolvers(srcName, nestedGraphqlObject, injectJoin)
			continue
		}

		if field.Resolve == nil {
			field.Resolve = a.resolveMiscField(srcName)
			graphqlObject.AddFieldConfig(field.Name, utils.GraphqlFieldDefinitionToField(field.Name, field))
		}
	}

	// Add an extra field to support on the fly joins
	graphqlObject.AddFieldConfig("_join", &graphql.Field{
		Type:        a.rootQueryType,
		Description: "Field injected to support on-the-fly joins",
		Resolve:     a.resolveJoin(),
	})
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

func getGraphqlTypeNameFromAST(t ast.Type) (typeName string, isList, isRequired bool) {
	switch v := t.(type) {
	case *ast.NonNull:
		typeName, isList, _ = getGraphqlTypeNameFromAST(v.Type)
		return typeName, isList, true
	case *ast.List:
		typeName, _, isRequired = getGraphqlTypeNameFromAST(v.Type)
		return typeName, true, isRequired
	case *ast.Named:
		return v.Name.Value, false, false
	}

	return "", false, false
}
