package schema

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/spaceuptech/space-cloud/config"
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
		log.Println("error ", err)
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
		case typeString:
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
