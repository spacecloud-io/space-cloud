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
func (m *Module) ParseSchema(crud config.Crud) error {
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

func (m *Module) schemaValidator(dbType, col string, fields map[string]interface{}) error {

	collectionFields := m.schema[dbType][col]
	for fieldKey, fieldValue := range collectionFields {
		// TODO: check if internal elements of list is required for this we need to change schmea parser
		// check if key is required
		value, ok := fields[fieldKey]
		if fieldValue.IsFieldTypeRequired {
			if !ok {
				return errors.New("Field " + fieldKey + " Not Present")
			}
		}

		// check type
		fieldType, err := m.checkType(value, fieldValue.Kind, dbType, fieldValue.TableJoin)
		if err != nil {
			return err
		}

		if fieldType != fieldValue.Kind {
			return errors.New("Wrong Type Wanter " + fieldValue.Kind + " got " + fieldType)
		}

	}

	return nil
}

// ValidateSchema checks data type
func (m *Module) ValidateSchema(dbType, col string, req *model.CreateRequest) error {

	v := make([]map[string]interface{}, 1)

	switch t := req.Document.(type) {
	case []map[string]interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	for _, fields := range v {
		err := m.schemaValidator(dbType, col, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) checkType(value interface{}, kind, dbType, col string) (string, error) {

	switch v := value.(type) {
	case int:
		// TODO: int64
		switch kind {
		case TypeDateTime:
			unitTimeInRFC3339 := time.Unix(int64(v), 0).Format(time.RFC3339)
			if unitTimeInRFC3339 == "" {
				return "", errors.New("Integer Wrong Date-Time Format")
			}
			return TypeDateTime, nil
		case TypeID:
			return TypeID, nil
		}
		return TypeInteger, nil

	case string:
		switch kind {
		case TypeDateTime:
			_, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return "", errors.New("String Wrong Date-Time Format")
			}
			return TypeDateTime, nil
		case TypeID:
			return TypeID, nil
		}
		return TypeString, nil

	case float32, float64:
		return TypeFloat, nil

	case bool:
		return TypeBoolean, nil

	case map[string]interface{}:
		err := m.schemaValidator(dbType, col, v)
		if err != nil {
			return "", err
		}
		return TypeJoin, nil

	case []interface{}:
		var str string
		var err error

		for _, value := range v {
			str, err = m.checkType(value, kind, dbType, col)
			if err != nil {
				return "", err
			}

			if str != kind {
				return "", errors.New("Wrong List Type Wanted " + kind + " got " + str)
			}
		}
		return str, nil

	}
	return "", errors.New("Checktype no match found")

}
