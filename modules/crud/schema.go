package crud

import (
	"errors"
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
	for dbName, v := range crud {
		collection := SchemaCollection{}
		for collectionName, v := range v.Collections {
			source := source.NewSource(&source.Source{
				Body: []byte(v.Schema),
			})
			// parse the source
			doc, err := parser.Parse(parser.ParseParams{Source: source})
			if err != nil {
				return nil, err
			}
			value, err := getCollectionSchema(doc, collectionName)
			if err != nil {
				return nil, err
			}
			collection[collectionName] = value
		}
		schema[dbName] = collection
	}
	return schema, nil
}

func getCollectionSchema(doc *ast.Document, collectionName string) (SchemaField, error) {
	fieldMap := SchemaField{}
	for _, v := range doc.Definitions {
		colName := v.(*ast.ObjectDefinition).Name.Value

		if colName != strings.Title(collectionName) {
			continue
		}
		for _, ve := range v.(*ast.ObjectDefinition).Fields {

			fieldTypeStuct := SchemaFieldType{
				Directive: DirectiveProperties{
					Value: DirectiveArgs{},
				},
			}
			if len(ve.Directives) > 0 {
				val := ve.Directives[0]
				argValue := map[string]string{}

				for _, x := range val.Arguments {

					val, _ := (utils.ParseGraphqlValue(x.Value, nil))
					argValue[x.Name.Value] = val.(string) // direvtive field name & value name
				}

				fieldTypeStuct.Directive.Kind = val.Name.Value
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
			getFieldType(fieldType.(*ast.List).Type, fieldTypeStuct, doc)

		}
	case kinds.Named:
		{
			myType := fieldType.(*ast.Named).Name.Value
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
					fieldTypeStuct.Kind = TypeRelation
					fieldTypeStuct.TableJoin = strings.ToLower(myType[0:1]) + myType[1:]
					if fieldTypeStuct.Directive.Kind != "relation" {
						fieldTypeStuct.Kind = TypeJoin
						nestedSchemaField, err := getCollectionSchema(doc, myType)
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
