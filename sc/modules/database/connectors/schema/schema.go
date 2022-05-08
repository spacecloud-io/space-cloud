package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
)

// Validate validates provided doc object against it's schema
func Validate(ctx context.Context, dbAlias, dbType, col string, collectionFields model.FieldSchemas, doc map[string]interface{}) (map[string]interface{}, error) {
	for schemaKey := range doc {
		if _, p := collectionFields[schemaKey]; !p {
			return nil, errors.New("The field " + schemaKey + " is not present in schema of " + col)
		}
	}

	mutatedDoc := map[string]interface{}{}
	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]

		if fieldValue.IsLinked {
			if ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("cannot insert value for a linked field %s", fieldKey), nil, nil)
			}
			continue
		}

		if fieldValue.IsAutoIncrement {
			continue
		}

		if !ok && fieldValue.IsDefault {
			defaultStringValue, isString := fieldValue.Default.(string)
			if fieldValue.Kind == model.TypeJSON && isString {
				var v interface{}
				_ = json.Unmarshal([]byte(defaultStringValue), &v)
				value = v
			} else {
				value = fieldValue.Default
			}
			ok = true
		}

		if fieldValue.IsCreatedAt || fieldValue.IsUpdatedAt {
			mutatedDoc[fieldKey] = time.Now().UTC()
			continue
		}

		if fieldValue.IsFieldTypeRequired {
			if fieldValue.Kind == model.TypeID && !ok {
				value = ksuid.New().String()
			} else if !ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("required field (%s) from table/colection (%s) not present in request", fieldKey, col), nil, nil)
			}
		}

		// check type
		val, err := checkType(ctx, dbAlias, dbType, col, value, fieldValue)
		if err != nil {
			return nil, err
		}

		mutatedDoc[fieldKey] = val
	}
	return mutatedDoc, nil
}

// Parser function parses the schema im module
func Parser(dbSchemas config.DatabaseSchemas) (model.CollectionSchemas, error) {
	schema := make(model.CollectionSchemas)
	for _, dbSchema := range dbSchemas {
		// Skip if no schema is provided
		if dbSchema.Schema == "" {
			continue
		}

		// Parse the graphql source
		s := source.NewSource(&source.Source{
			Body: []byte(dbSchema.Schema),
		})
		doc, err := parser.Parse(parser.ParseParams{Source: s})
		if err != nil {
			return nil, err
		}

		// Parse graphql ast to map of fields
		value, err := getCollectionSchema(doc, dbSchema.DbAlias, dbSchema.Table)
		if err != nil {
			return nil, err
		}

		// We will skip this table if it has only one field or less
		if len(value) <= 1 { // schema might have an id by default
			continue
		}

		schema[dbSchema.Table] = value
	}
	return schema, nil
}

// GetConstraintName generates constraint name for joint fields
func GetConstraintName(tableName, columnName string) string {
	return fmt.Sprintf("c_%s_%s", tableName, columnName)
}

// AddInternalLinks links the internal links based on the foreign keys
func AddInternalLinks(dbAlias string, schemas model.CollectionSchemas) {
	for tableName, fieldSchemas := range schemas {
		for _, fieldSchema := range fieldSchemas {
			if fieldSchema.IsForeign {
				// Add a linked field in this table
				fieldSchemas[fieldSchema.JointTable.Table] = &model.FieldType{
					FieldName: fieldSchema.JointTable.Table,
					IsLinked:  true,
					IsList:    false,
					Kind:      fieldSchema.JointTable.Table,
					LinkedTable: &model.LinkProperties{
						DB:    dbAlias,
						Table: fieldSchema.JointTable.Table,
						From:  fieldSchema.FieldName,
						To:    fieldSchema.JointTable.To,
					},
				}

				// Add a linked field in the joint table
				jointTableFieldsSchemas, p := schemas[fieldSchema.JointTable.Table]
				if !p {
					continue
				}

				jointTableFieldsSchemas[tableName] = &model.FieldType{
					FieldName: tableName,
					IsLinked:  true,
					IsList:    true,
					Kind:      tableName,
					LinkedTable: &model.LinkProperties{
						DB:    dbAlias,
						Table: tableName,
						From:  fieldSchema.JointTable.To,
						To:    fieldSchema.FieldName,
					},
				}

				// TODO: Add support for automatic many to many joins
			}
		}
	}
}
