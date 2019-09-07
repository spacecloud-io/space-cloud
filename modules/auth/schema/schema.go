package schema

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
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// TODO: check graphql types
type Schema struct {
	SchemaDoc schemaType
	crud      *crud.Module
	project   string
}

// Init creates a new instance of the schema object
func Init(crud *crud.Module) *Schema {
	return &Schema{SchemaDoc: schemaType{}, crud: crud}
}

func (s *Schema) SetProject(project string) {
	s.project = project
}

// ParseSchema Initializes Schema field in Module struct
func (s *Schema) ParseSchema(crud config.Crud) error {

	schema := make(schemaType, len(crud))
	for dbName, v := range crud {
		collection := schemaCollection{}
		for collectionName, v := range v.Collections {
			if v.Schema == "" {
				continue
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
	s.SchemaDoc = schema
	return nil
}

func getCollectionSchema(doc *ast.Document, collectionName string) (schemaField, error) {
	fieldMap := schemaField{}
	for _, v := range doc.Definitions {
		colName := v.(*ast.ObjectDefinition).Name.Value

		if colName != strings.Title(collectionName) {
			continue
		}
		for _, ve := range v.(*ast.ObjectDefinition).Fields {

			fieldTypeStuct := schemaFieldType{
				directive: directiveProperties{
					Value: directiveArgs{},
				},
			}
			if len(ve.Directives) > 0 {
				val := ve.Directives[0]
				argValue := map[string]string{}

				for _, x := range val.Arguments {

					val, _ := (utils.ParseGraphqlValue(x.Value, nil))
					argValue[x.Name.Value] = val.(string) // direvtive field name & value name
				}

				fieldTypeStuct.directive.Kind = val.Name.Value
				fieldTypeStuct.directive.Value = argValue
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

func getFieldType(fieldType ast.Type, fieldTypeStuct *schemaFieldType, doc *ast.Document) error {
	switch fieldType.GetKind() {
	case kinds.NonNull:
		{
			fieldTypeStuct.isFieldTypeRequired = true
			getFieldType(fieldType.(*ast.NonNull).Type, fieldTypeStuct, doc)
		}
	case kinds.List:
		{
			fieldTypeStuct.isList = true
			getFieldType(fieldType.(*ast.List).Type, fieldTypeStuct, doc)

		}
	case kinds.Named:
		{
			myType := fieldType.(*ast.Named).Name.Value
			switch myType {
			case typeString, typeEnum:
				fieldTypeStuct.Kind = typeString
			case typeID:
				fieldTypeStuct.Kind = typeID
			case typeDateTime:
				fieldTypeStuct.Kind = typeDateTime
			case typeFloat:
				fieldTypeStuct.Kind = typeFloat
			case typeInteger:
				fieldTypeStuct.Kind = typeInteger
			case typeBoolean:
				fieldTypeStuct.Kind = typeBoolean
			case typeJSON:
				fieldTypeStuct.Kind = typeJSON
			default:
				{
					fieldTypeStuct.Kind = typeRelation
					fieldTypeStuct.tableJoin = strings.ToLower(myType[0:1]) + myType[1:]
					if fieldTypeStuct.directive.Kind != "relation" {
						fieldTypeStuct.Kind = typeJoin
						nestedschemaField, err := getCollectionSchema(doc, myType)
						if err != nil {
							return err
						}
						fieldTypeStuct.nestedObject = nestedschemaField
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

func (s *Schema) schemaValidator(collectionFields schemaField, doc map[string]interface{}) (map[string]interface{}, error) {
	if len(collectionFields) == 0 {
		return doc, nil
	}

	mutatedDoc := map[string]interface{}{}

	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]
		if fieldValue.isFieldTypeRequired {
			if !ok {
				return nil, errors.New("Field " + fieldKey + " Not Present")
			}
		}

		// check type
		val, err := s.checkType(value, fieldValue)
		if err != nil {
			return nil, err
		}
		mutatedDoc[fieldKey] = val
	}

	return mutatedDoc, nil
}

// ValidateCreateOperation
func (s *Schema) ValidateCreateOperation(dbType, col string, req *model.CreateRequest) error {

	if s.SchemaDoc == nil {
		return errors.New("Schema not initialized")
	}

	v := make([]interface{}, 0)

	switch t := req.Document.(type) {
	case []interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	collection, ok := s.SchemaDoc[dbType]
	if !ok {
		return errors.New("No db was found named " + dbType)
	}
	collectionFields, ok := collection[col]
	if !ok {
		return errors.New("No collection or table was found named " + col)
	}

	for index, doc := range v {
		newDoc, err := s.schemaValidator(collectionFields, doc.(map[string]interface{}))
		if err != nil {
			return err
		}
		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}

func (s *Schema) checkType(value interface{}, fieldValue *schemaFieldType) (interface{}, error) {

	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case typeDateTime:
			return time.Unix(int64(v), 0), nil
		case typeID, typeInteger:
			return value, nil
		default:
			return nil, errors.New("Integer wrong type wanted " + fieldValue.Kind + " got Integer")
		}

	case string:
		switch fieldValue.Kind {
		case typeDateTime:
			unitTimeInRFC3339, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, errors.New("String Wrong Date-Time Format")
			}
			return unitTimeInRFC3339, nil
		case typeID, typeString:
			return value, nil
		default:
			return nil, errors.New("String wrong type wanted " + fieldValue.Kind + " got String")
		}

	case float32, float64:
		switch fieldValue.Kind {
		case typeFloat:
			return value, nil
		default:
			return nil, errors.New("Float wrong type wanted " + fieldValue.Kind + " got Float")
		}
	case bool:
		switch fieldValue.Kind {
		case typeBoolean:
			return value, nil
		default:
			return nil, errors.New("Bool wrong type wanted " + fieldValue.Kind + " got Bool")
		}
	case map[string]interface{}:
		return s.schemaValidator(fieldValue.nestedObject, v)

	case []interface{}:
		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := s.checkType(value, fieldValue)
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
