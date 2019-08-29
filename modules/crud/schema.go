package crud

import (
	"errors"
	"strings"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// TODO: check graphql types

// ParseSchema Initializes Schema field in Module struct
func (m *Module) parseSchema(crud config.Crud) error {

	schema := make(SchemaType, len(crud))
	for dbName, v := range crud {
		collection := SchemaCollection{}
		for collectionName, v := range v.Collections {
			if v.Schema == "" {
				return errors.New("Schema is empty")
			}
			source := source.NewSource(&source.Source{
				Body: []byte(v.Schema),
			})
			// parse the source
			doc, err := parser.Parse(parser.ParseParams{Source: source})
			if err != nil {
				return err
			}
			value, err := getCollectionSchema(doc, collectionName)
			if err != nil {
				return err
			}
			collection[collectionName] = value
		}
		schema[dbName] = collection
	}
	m.schema = schema
	return nil
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

func (m *Module) schemaValidator(collectionFields SchemaField, doc map[string]interface{}) (map[string]interface{}, error) {
	mutatedDoc := map[string]interface{}{}

	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]
		if fieldValue.IsFieldTypeRequired {
			if !ok {
				return nil, errors.New("Field " + fieldKey + " Not Present")
			}
		}

		// check type
		val, err := m.checkType(value, fieldValue)
		if err != nil {
			return nil, err
		}
		mutatedDoc[fieldKey] = val
	}

	return mutatedDoc, nil
}

// ValidateSchema checks data type
func (m *Module) ValidateSchema(dbType, col string, req *model.CreateRequest) error {

	if m.schema == nil {
		return errors.New("Schema not initialized")
	}

	v := make([]map[string]interface{}, 0)

	switch t := req.Document.(type) {
	case []map[string]interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	collection, ok := m.schema[dbType]
	if !ok {
		return errors.New("No db was found named " + dbType)
	}
	collectionFields, ok := collection[col]
	if !ok {
		return errors.New("No collection or table was found named " + col)
	}

	for index, fields := range v {

		newDoc, err := m.schemaValidator(collectionFields, fields)
		if err != nil {
			return err
		}
		v[index] = newDoc
	}

	return nil
}

func (m *Module) checkType(value interface{}, fieldValue *SchemaFieldType) (interface{}, error) {

	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case TypeDateTime:
			unitTimeInRFC3339 := time.Unix(int64(v), 0).Format(time.RFC3339)
			if unitTimeInRFC3339 == "" {
				return nil, errors.New("Integer Wrong Date-Time Format")
			}
			return unitTimeInRFC3339, nil
		case TypeID, TypeInteger:
			return value, nil
		default:
			return nil, errors.New("Integer wrong type wanted " + fieldValue.Kind + " got Integer")
		}

	case string:
		switch fieldValue.Kind {
		case TypeDateTime:
			unitTimeInRFC3339, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, errors.New("String Wrong Date-Time Format")
			}
			return unitTimeInRFC3339, nil
		case TypeID, TypeString:
			return value, nil
		default:
			return nil, errors.New("String wrong type wanted " + fieldValue.Kind + " got String")
		}

	case float32, float64:
		switch fieldValue.Kind {
		case TypeFloat:
			return value, nil
		default:
			return nil, errors.New("Float wrong type wanted " + fieldValue.Kind + " got Float")
		}
	case bool:
		switch fieldValue.Kind {
		case TypeBoolean:
			return value, nil
		default:
			return nil, errors.New("Bool wrong type wanted " + fieldValue.Kind + " got Bool")
		}
	case map[string]interface{}:
		newDoc, err := m.schemaValidator(fieldValue.NestedObject, v)
		if err != nil {
			return nil, err
		}
		return newDoc, nil

	case []interface{}:
		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := m.checkType(value, fieldValue)
			if err != nil {
				return nil, err
			}
			arr[index] = val
		}
		return arr, nil
	default:
		return nil, errors.New("No matching type found")
	}
}
