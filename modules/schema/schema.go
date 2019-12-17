package schema

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Schema data stucture for schema package
type Schema struct {
	lock               sync.RWMutex
	SchemaDoc          schemaType
	crud               *crud.Module
	project            string
	config             config.Crud
	removeProjectScope bool
}

// Init creates a new instance of the schema object
func Init(crud *crud.Module, removeProjectScope bool) *Schema {
	return &Schema{SchemaDoc: schemaType{}, crud: crud, removeProjectScope: removeProjectScope}
}

// SetConfig modifies the tables according to the schema on save
func (s *Schema) SetConfig(conf config.Crud, project string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = conf
	s.project = project

	if err := s.parseSchema(conf); err != nil {
		return err
	}

	return nil
}

func (s *Schema) GetSchema(dbType, col string) (SchemaFields, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	dbSchema, p := s.SchemaDoc[dbType]
	if !p {
		return nil, false
	}

	colSchema, p := dbSchema[col]
	if !p {
		return nil, false
	}

	fields := make(SchemaFields, len(colSchema))
	for k, v := range colSchema {
		fields[k] = v
	}

	return fields, true
}

// parseSchema Initializes Schema field in Module struct
func (s *Schema) parseSchema(crud config.Crud) error {

	schema, err := s.parser(crud)
	if err != nil {
		return err
	}
	s.SchemaDoc = schema
	return nil
}

func (s *Schema) parser(crud config.Crud) (schemaType, error) {

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
				return nil, err
			}
			value, err := getCollectionSchema(doc, dbName, collectionName)
			if err != nil {
				return nil, err
			}

			if len(value) <= 1 { // schema might have an id by default
				continue
			}
			collection[strings.ToLower(collectionName[0:1])+collectionName[1:]] = value
		}
		schema[dbName] = collection
	}
	return schema, nil
}

func getCollectionSchema(doc *ast.Document, dbName, collectionName string) (SchemaFields, error) {
	var isCollectionFound bool

	fieldMap := SchemaFields{}
	for _, v := range doc.Definitions {
		colName := v.(*ast.ObjectDefinition).Name.Value

		if colName != collectionName {
			continue
		}

		// Mark the collection as found
		isCollectionFound = true

		for _, field := range v.(*ast.ObjectDefinition).Fields {

			fieldTypeStuct := SchemaFieldType{
				FieldName: field.Name.Value,
			}
			if len(field.Directives) > 0 {
				// Loop over all directives

				for _, directive := range field.Directives {
					switch directive.Name.Value {
					case directivePrimary:
						fieldTypeStuct.IsPrimary = true
					case directiveUnique:
						fieldTypeStuct.IsUnique = true
					case directiveCreatedAt:
						fieldTypeStuct.IsCreatedAt = true
					case directiveUpdatedAt:
						fieldTypeStuct.IsUpdatedAt = true
					case directiveDefault:
						fieldTypeStuct.IsDefault = true

						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "value":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)

								fieldTypeStuct.Default = val
							}
							if value == "" {
								return nil, errors.New("The default value can not be null")
							}
						}
					case directiveLink:
						fieldTypeStuct.IsLinked = true
						fieldTypeStuct.LinkedTable = &TableProperties{DBType: dbName}
						kind, err := getFieldType(dbName, field.Type, &fieldTypeStuct, doc)
						if err != nil {
							return nil, err
						}
						fieldTypeStuct.LinkedTable.Table = kind

						// Load the from and to fields. If either is not available, we will throw an error.
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "table":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.Table = val.(string)
							case "from":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.From = val.(string)
							case "to":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.To = val.(string)
							case "field":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.Field = val.(string)
							}
						}

						// Throw an error if from and to are unavailable
						if fieldTypeStuct.LinkedTable.From == "" || fieldTypeStuct.LinkedTable.To == "" {
							return nil, errors.New("link directive must be accompanied with to and from fields")
						}

					case directiveForeign:
						fieldTypeStuct.IsForeign = true
						fieldTypeStuct.JointTable = &TableProperties{}
						fieldTypeStuct.JointTable.Table = strings.Split(field.Name.Value, "_")[0]
						fieldTypeStuct.JointTable.From = field.Name.Value
						fieldTypeStuct.JointTable.To = "id"

						// Load the joint table name and field
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "table":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.Table = val.(string)

							case "field":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.To = val.(string)
							}
						}
					}
				}
			}

			kind, err := getFieldType(dbName, field.Type, &fieldTypeStuct, doc)
			if err != nil {
				return nil, err
			}
			fieldTypeStuct.Kind = kind
			fieldMap[field.Name.Value] = &fieldTypeStuct
		}
	}

	// Throw an error if the collection wasn't found
	if !isCollectionFound {
		return nil, fmt.Errorf("collection %s could not be found in schema", collectionName)
	}
	return fieldMap, nil
}

func getFieldType(dbName string, fieldType ast.Type, fieldTypeStuct *SchemaFieldType, doc *ast.Document) (string, error) {
	switch fieldType.GetKind() {
	case kinds.NonNull:
		fieldTypeStuct.IsFieldTypeRequired = true
		return getFieldType(dbName, fieldType.(*ast.NonNull).Type, fieldTypeStuct, doc)
	case kinds.List:
		// Lists are not allowed for primary and foreign keys
		if fieldTypeStuct.IsPrimary || fieldTypeStuct.IsForeign {
			return "", fmt.Errorf("invalid type for field %s - primary and foreign keys cannot be made on lists", fieldTypeStuct.FieldName)
		}

		fieldTypeStuct.IsList = true
		return getFieldType(dbName, fieldType.(*ast.List).Type, fieldTypeStuct, doc)

	case kinds.Named:
		myType := fieldType.(*ast.Named).Name.Value
		switch myType {
		case typeString, typeEnum:
			return typeString, nil
		case TypeID:
			return TypeID, nil
		case typeDateTime:
			return typeDateTime, nil
		case typeFloat:
			return typeFloat, nil
		case typeInteger:
			return typeInteger, nil
		case typeBoolean:
			return typeBoolean, nil
		case typeJSON:
			return typeJSON, nil
		default:
			if fieldTypeStuct.IsLinked {
				// Since the field is actually a link. We'll store the type as is. This type must correspond to a table or a primitive type
				// or else the link won't work. It's upto the user to make sure of that.
				return myType, nil
			}

			// The field is a nested type. Update the nestedObject field and return typeObject. This is a side effect.
			nestedschemaField, err := getCollectionSchema(doc, dbName, myType)
			if err != nil {
				return "", err
			}
			fieldTypeStuct.nestedObject = nestedschemaField

			return typeObject, nil
		}
	default:
		return "", fmt.Errorf("invalid field kind `%s` provided for field `%s`", fieldType.GetKind(), fieldTypeStuct.FieldName)
	}
}

func (s *Schema) schemaValidator(col string, collectionFields SchemaFields, doc map[string]interface{}) (map[string]interface{}, error) {

	mutatedDoc := map[string]interface{}{}

	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]

		if fieldValue.IsLinked {
			if ok {
				return nil, fmt.Errorf("cannot insert value for a linked field %s", fieldKey)
			}

			continue
		}

		if !ok && fieldValue.IsDefault {
			mutatedDoc[fieldKey] = fieldValue.Default
			continue
		}

		if fieldValue.Kind == TypeID && !ok {
			value = ksuid.New().String()
			ok = true
		}

		if fieldValue.IsCreatedAt || fieldValue.IsUpdatedAt {
			mutatedDoc[fieldKey] = time.Now().UTC()
			continue
		}

		if fieldValue.IsFieldTypeRequired {
			if !ok {
				return nil, fmt.Errorf("required field %s from collection %s not present in request", fieldKey, col)
			}
		}

		// check type
		val, err := s.checkType(col, value, fieldValue)
		if err != nil {
			return nil, err
		}

		mutatedDoc[fieldKey] = val
	}

	return mutatedDoc, nil
}

// ValidateCreateOperation validates schema on create operation
func (s *Schema) ValidateCreateOperation(dbType, col string, req *model.CreateRequest) error {

	if s.SchemaDoc == nil {
		return errors.New("schema not initialized")
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
		return nil
	}

	for index, doc := range v {
		newDoc, err := s.schemaValidator(col, collectionFields, doc.(map[string]interface{}))
		if err != nil {
			return err
		}

		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}

func (s *Schema) checkType(col string, value interface{}, fieldValue *SchemaFieldType) (interface{}, error) {

	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case typeDateTime:
			return time.Unix(int64(v)/1000, 0), nil
		case typeInteger:
			return value, nil
		case typeFloat:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case string:
		switch fieldValue.Kind {
		case typeDateTime:
			unitTimeInRFC3339, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("invalid datetime format recieved for field %s in collection %s - use RFC3339 fromat", fieldValue.FieldName, col)
			}
			return unitTimeInRFC3339, nil
		case TypeID, typeString:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got String", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case float32, float64:
		switch fieldValue.Kind {
		case typeDateTime:
			return time.Unix(int64(v.(float64))/1000, 0), nil
		case typeFloat:
			return value, nil
		case typeInteger:
			return int64(value.(float64)), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Float", fieldValue.FieldName, col, fieldValue.Kind)
		}
	case bool:
		switch fieldValue.Kind {
		case typeBoolean:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Bool", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case map[string]interface{}:
		// TODO: allow this operation for nested insert using links
		if fieldValue.Kind != typeObject {
			return nil, fmt.Errorf("invalid type received for field %s in collection %s", fieldValue.FieldName, col)
		}

		return s.schemaValidator(col, fieldValue.nestedObject, v)

	case []interface{}:
		// TODO: allow this operation for nested insert using links
		if fieldValue.Kind != typeObject {
			return nil, fmt.Errorf("invalid type received for field %s in collection %s", fieldValue.FieldName, col)
		}

		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := s.checkType(col, value, fieldValue)
			if err != nil {
				return nil, err
			}
			arr[index] = val
		}
		return arr, nil
	default:
		if !fieldValue.IsFieldTypeRequired || fieldValue.IsDefault {
			return nil, nil
		}

		return nil, fmt.Errorf("no matching type found for field %s in collection %s", fieldValue.FieldName, col)
	}
}
