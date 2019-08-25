package crud

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// ParseSchema Initializes Schema field in Module struct
func (m *Module) ParseSchema(crud config.Crud) (SchemaType, error) {
	schema := make(SchemaType, len(crud))
	collection := SchemaCollection{}
	for dbName, v := range crud {
		for collectionName, v := range v.Collections {
			source := source.NewSource(&source.Source{
				Body: []byte(v.Schema),
			})
			// parse the source
			doc, err := parser.Parse(parser.ParseParams{Source: source})
			if err != nil {
				return nil, err
			}
			value, _ := getSchemaDetails(doc, collectionName)
			collection[collectionName] = value
		}
		schema[dbName] = collection
	}
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
	return schema, nil
}

func getSchemaDetails(doc *ast.Document, collectionName string) (SchemaField, error) {
	fieldMap := SchemaField{}
	for _, v := range doc.Definitions {
		for _, ve := range v.(*ast.ObjectDefinition).Fields {

			colName := v.(*ast.ObjectDefinition).Name.Value
			// fmt.Println("Collection Name", collectionName, colName)

			if colName != strings.Title(collectionName) {
				continue
			}
			fieldTypeStuct := SchemaFieldType{
				Directive: DirectiveProperties{
					Value: DirectiveArgs{},
				},
			}

			for _, w := range ve.Directives {
				argValue := map[string]string{}
				for _, x := range w.Arguments {

					val, _ := (utils.ParseGraphqlValue(x.Value, nil))
					argValue[x.Name.Value] = val.(string) // direvtive field name & value name
				}
				fieldTypeStuct.Directive.Kind = w.Name.Value
				fieldTypeStuct.Directive.Value = argValue
			}
			err := getFieldType(ve.Type, &fieldTypeStuct, doc)
			if err != nil {
				return nil, err
			}
			fieldMap[ve.Name.Value] = &fieldTypeStuct
		}
	}

	return fieldMap, nil
}

func getFieldType(fieldType ast.Type, fieldTypeStuct *SchemaFieldType, doc *ast.Document) error {
	switch fieldType.GetKind() {
	case kinds.NonNull:
		{
			fieldTypeStuct.IsFieldTypeRequired = true
			getFieldType(fieldType.(*ast.NonNull).Type, fieldTypeStuct, doc)

		}
	case kinds.List:
		{
			fieldTypeStuct.IsList = true
			if fieldType.(*ast.List).Type.GetKind() == kinds.Named {
				getFieldType(fieldType.(*ast.List).Type, fieldTypeStuct, doc)
			} else {
				getFieldType(fieldType.(*ast.List).Type.(*ast.NonNull).Type, fieldTypeStuct, doc)
			}
		}
	case kinds.Named:
		{
			myType := fieldType.(*ast.Named).Name.Value
			fmt.Println("Type:", myType)
			switch myType {
			case TypeString, TypeEnum:
				fieldTypeStuct.Kind = TypeString
			case TypeID:
				fieldTypeStuct.Kind = TypeID
			case TypeDateTime:
				fieldTypeStuct.Kind = TypeDateTime
			case TypeFloat:
				fieldTypeStuct.Kind = TypeFloat
			case TypeInteger:
				fieldTypeStuct.Kind = TypeInteger
			case TypeBoolean:
				fieldTypeStuct.Kind = TypeBoolean
			case TypeJSON:
				fieldTypeStuct.Kind = TypeJSON
			default:
				{
					fmt.Println("Mytype here", myType)
					fieldTypeStuct.Kind = TypeRelation
					fieldTypeStuct.TableJoin = strings.ToLower(myType[0:1]) + myType[1:]
					if fieldTypeStuct.Directive.Kind != "relation" {
						fieldTypeStuct.Kind = TypeJoin
						nestedSchemaField, err := getSchemaDetails(doc, myType)
						if err != nil {
							return err
						}
						fieldTypeStuct.NestedObject = nestedSchemaField
					}

				}
			}
		}
	default:
		{
			return errors.New("Wrong Field Type")
		}
	}
	return nil
}

// ValidateSchema checks data type
func (m *Module) ValidateSchema() {

}
