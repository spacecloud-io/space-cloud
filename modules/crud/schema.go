package crud

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// ParseSchema Initializes Schema field in Module struct
func (m *Module) ParseSchema(crud config.Crud) ([]SchemaType, error) {
	schemaSlice := make([]SchemaType, len(crud))
	for dbName, v := range crud {
		for _, v := range v.Collections {
			source := source.NewSource(&source.Source{
				Body: []byte(v.Schema),
				Name: "GraphQL request",
			})
			// parse the source
			doc, err := parser.Parse(parser.ParseParams{Source: source})
			if err != nil {
				return nil, err
			}
			value := getSchemaDetails(doc, dbName)
			// fmt.Println(value)
			// for _, v := range value {
			// 	for _, v := range v {
			// 		for _, v := range v {
			// 			fmt.Println("For ", v)
			// 		}
			// 	}
			// }
			schemaSlice = append(schemaSlice, value)
		}
	}
	b, err := json.MarshalIndent(schemaSlice, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
	return schemaSlice, nil
}

func getFieldType(ve ast.Type, fieldTypeStuct *schemaFieldType) {

	if ve.GetKind() == kinds.NonNull {
		myType := ve.(*ast.NonNull)
		fieldTypeStuct.IsFieldTypeRequired = true
		ve = myType.Type
		if myType.Type.String() == kinds.List {
			fieldTypeStuct.IsList = true
			//getFieldType(myType.Type.(*ast.List).Type, fieldTypeStuct)
			ve = myType.Type.(*ast.List).Type.(*ast.NonNull).Type
		}
	}

	if ve.GetKind() == kinds.Named {
		myType := ve.(*ast.Named)
		fmt.Println("Type:", myType.Name.Value)

		switch myType.Name.Value {
		case "String":
			fieldTypeStuct.Kind = TypeString
		case "ID":
			fieldTypeStuct.Kind = TypeID
		case "DateTime":
			fieldTypeStuct.Kind = TypeDateTime
		case "Float":
			fieldTypeStuct.Kind = TypeFloat
		case "Integer":
			fieldTypeStuct.Kind = TypeInteger
		case "Boolean":
			fieldTypeStuct.Kind = TypeBoolean
		case "Enum":
			fieldTypeStuct.Kind = TypeEnum
		case "Json":
			fieldTypeStuct.Kind = TypeJSON
		default:
			fieldTypeStuct.Kind = TypeJoin
		}
	}
}

func getSchemaDetails(doc *ast.Document, dbName string) SchemaType {
	var m = make(SchemaType)
	collectionMap := schemaCollection{}
	for _, v := range doc.Definitions {
		fieldMap := schemaField{}
		for _, ve := range v.(*ast.ObjectDefinition).Fields {
			fieldTypeStuct := schemaFieldType{
				DirectiveType: map[string]schemaFieldDirectiveArgs{},
			}
			getFieldType(ve.Type, &fieldTypeStuct)
			for _, w := range ve.Directives {
				fieldTypeStuct.DirectiveType[w.Name.Value] = nil // directive name
				argValue := map[string]string{}
				for _, x := range w.Arguments {

					val, _ := (utils.ParseGraphqlValue(x.Value, nil))
					argValue[x.Name.Value] = val.(string) // direvtive field name & value name
				}
				fieldTypeStuct.DirectiveType[w.Name.Value] = argValue
			}
			fieldMap[ve.Name.Value] = &fieldTypeStuct
		}
		collectionMap[v.(*ast.ObjectDefinition).Name.Value] = fieldMap
	}
	m[dbName] = collectionMap // db name
	return m
}

// ValidateSchema checks data type
func (m *Module) ValidateSchema() {

}
