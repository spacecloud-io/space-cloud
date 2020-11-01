package schema

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Schema data stucture for schema package
type Schema struct {
	lock      sync.RWMutex
	SchemaDoc model.Type
	crud      model.CrudSchemaInterface
	project   string
	dbSchemas config.DatabaseSchemas
	clusterID string
}

// Init creates a new instance of the schema object
func Init(clusterID string, crud model.CrudSchemaInterface) *Schema {
	return &Schema{clusterID: clusterID, SchemaDoc: model.Type{}, crud: crud}
}

// SetConfig modifies the tables according to the schema on save
func (s *Schema) SetConfig(c config.DatabaseSchemas, project string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.dbSchemas = c
	s.project = project

	if err := s.parseSchema(c); err != nil {
		return err
	}

	return nil
}

// GetSchema function gets schema
func (s *Schema) GetSchema(dbAlias, col string) (model.Fields, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	dbSchema, p := s.SchemaDoc[dbAlias]
	if !p {
		return nil, false
	}

	colSchema, p := dbSchema[col]
	if !p {
		return nil, false
	}

	fields := make(model.Fields, len(colSchema))
	for k, v := range colSchema {
		fields[k] = v
	}

	return fields, true
}

// parseSchema Initializes Schema field in Module struct
func (s *Schema) parseSchema(crud config.DatabaseSchemas) error {
	schema, err := s.Parser(crud)
	if err != nil {
		return err
	}
	s.SchemaDoc = schema
	return nil
}

// Parser function parses the schema im module
func (s *Schema) Parser(dbSchemas config.DatabaseSchemas) (model.Type, error) {
	schema := make(model.Type)
	for _, dbSchema := range dbSchemas {
		if dbSchema.Schema == "" {
			continue
		}
		s := source.NewSource(&source.Source{
			Body: []byte(dbSchema.Schema),
		})
		// parse the source
		doc, err := parser.Parse(parser.ParseParams{Source: s})
		if err != nil {
			return nil, err
		}
		value, err := getCollectionSchema(doc, dbSchema.DbAlias, dbSchema.Table)
		if err != nil {
			return nil, err
		}

		if len(value) <= 1 { // schema might have an id by default
			continue
		}
		_, ok := schema[dbSchema.DbAlias]
		if !ok {
			schema[dbSchema.DbAlias] = model.Collection{dbSchema.Table: value}
		} else {
			schema[dbSchema.DbAlias][dbSchema.Table] = value
		}
	}
	return schema, nil
}

func getCollectionSchema(doc *ast.Document, dbName, collectionName string) (model.Fields, error) {
	var isCollectionFound bool

	fieldMap := model.Fields{}
	for _, v := range doc.Definitions {
		colName := v.(*ast.ObjectDefinition).Name.Value
		if colName != collectionName {
			continue
		}

		// Mark the collection as found
		isCollectionFound = true

		for _, field := range v.(*ast.ObjectDefinition).Fields {

			if field.Type == nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Type not provided for collection/table schema (%s) with field (%s)", collectionName, field.Name.Value), nil, nil)
			}

			fieldTypeStuct := model.FieldType{
				FieldName:  field.Name.Value,
				TypeIDSize: 50,
			}
			if len(field.Directives) > 0 {
				// Loop over all directives

				for _, directive := range field.Directives {
					switch directive.Name.Value {
					case model.DirectivePrimary:
						fieldTypeStuct.IsPrimary = true
					case model.DirectiveCreatedAt:
						fieldTypeStuct.IsCreatedAt = true
					case model.DirectiveUpdatedAt:
						fieldTypeStuct.IsUpdatedAt = true
					case model.DirectiveVarcharSize:
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "value":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								size, ok := val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								fieldTypeStuct.TypeIDSize = size
							}
						}
					case model.DirectiveDefault:
						fieldTypeStuct.IsDefault = true

						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "value":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.Default = val
							}
						}
						if fieldTypeStuct.Default == nil {
							return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Default directive must be accompanied with value field", nil, nil)
						}
					case model.DirectiveIndex, model.DirectiveUnique:
						fieldTypeStuct.IsIndex = true
						fieldTypeStuct.IsUnique = directive.Name.Value == model.DirectiveUnique
						fieldTypeStuct.IndexInfo = &model.TableProperties{Group: fieldTypeStuct.FieldName, Order: model.DefaultIndexOrder, Sort: model.DefaultIndexSort}
						for _, arg := range directive.Arguments {
							var ok bool
							switch arg.Name.Value {
							case "name", "group":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.IndexInfo.Group, ok = val.(string)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
							case "order":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.IndexInfo.Order, ok = val.(int)
								if !ok {
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected int", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
							case "sort":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								sort, ok := val.(string)
								if !ok || (sort != "asc" && sort != "desc") {
									if !ok {
										return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unexpected argument type provided for field (%s) directinve @(%s) argument (%s) got (%v) expected string", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
									}
									return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unknow value provided for field (%s) directinve @(%s) argument (%s) got (%v) expected either (asc) or (desc)", fieldTypeStuct.FieldName, directive.Name.Value, arg.Name.Value, reflect.TypeOf(val)), nil, map[string]interface{}{"arg": arg.Name.Value})
								}
								fieldTypeStuct.IndexInfo.Sort = sort
							}
						}
					case model.DirectiveLink:
						fieldTypeStuct.IsLinked = true
						fieldTypeStuct.LinkedTable = &model.TableProperties{DBType: dbName}
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
							case "db":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.LinkedTable.DBType = val.(string)
							}
						}

						// Throw an error if from and to are unavailable
						if fieldTypeStuct.LinkedTable.From == "" || fieldTypeStuct.LinkedTable.To == "" {
							return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Link directive must be accompanied with (to) and (from) arguments", nil, nil)
						}

					case model.DirectiveForeign:
						fieldTypeStuct.IsForeign = true
						fieldTypeStuct.JointTable = &model.TableProperties{}
						fieldTypeStuct.JointTable.Table = strings.Split(field.Name.Value, "_")[0]
						fieldTypeStuct.JointTable.To = "id"
						fieldTypeStuct.JointTable.OnDelete = "NO ACTION"

						// Load the joint table name and field
						for _, arg := range directive.Arguments {
							switch arg.Name.Value {
							case "table":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.Table = val.(string)

							case "field", "to":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.To = val.(string)

							case "onDelete":
								val, _ := utils.ParseGraphqlValue(arg.Value, nil)
								fieldTypeStuct.JointTable.OnDelete = val.(string)
								if fieldTypeStuct.JointTable.OnDelete == "cascade" {
									fieldTypeStuct.JointTable.OnDelete = "CASCADE"
								} else {
									fieldTypeStuct.JointTable.OnDelete = "NO ACTION"
								}
							}
						}
						fieldTypeStuct.JointTable.ConstraintName = getConstraintName(collectionName, fieldTypeStuct.FieldName)
					default:
						return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unknown directive (%s) provided for field (%s)", directive.Name.Value, fieldTypeStuct.FieldName), nil, nil)
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
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Collection/Table (%s) not found in schema", collectionName), nil, nil)
	}
	return fieldMap, nil
}

func getFieldType(dbName string, fieldType ast.Type, fieldTypeStuct *model.FieldType, doc *ast.Document) (string, error) {
	switch fieldType.GetKind() {
	case kinds.NonNull:
		fieldTypeStuct.IsFieldTypeRequired = true
		return getFieldType(dbName, fieldType.(*ast.NonNull).Type, fieldTypeStuct, doc)
	case kinds.List:
		// Lists are not allowed for primary and foreign keys
		if fieldTypeStuct.IsPrimary || fieldTypeStuct.IsForeign {
			return "", helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Primary and foreign keys directives cannot be added on field (%s) with type lists", fieldTypeStuct.FieldName), nil, nil)
		}

		fieldTypeStuct.IsList = true
		return getFieldType(dbName, fieldType.(*ast.List).Type, fieldTypeStuct, doc)

	case kinds.Named:
		myType := fieldType.(*ast.Named).Name.Value
		switch myType {
		case model.TypeString, model.TypeEnum:
			return model.TypeString, nil
		case model.TypeID:
			return model.TypeID, nil
		case model.TypeDateTime:
			return model.TypeDateTime, nil
		case model.TypeFloat:
			return model.TypeFloat, nil
		case model.TypeInteger:
			return model.TypeInteger, nil
		case model.TypeBoolean:
			return model.TypeBoolean, nil
		case model.TypeJSON:
			return model.TypeJSON, nil
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
			fieldTypeStuct.NestedObject = nestedschemaField

			return model.TypeObject, nil
		}
	default:
		return "", helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid field kind `%s` provided for field `%s`", fieldType.GetKind(), fieldTypeStuct.FieldName), nil, nil)
	}
}
