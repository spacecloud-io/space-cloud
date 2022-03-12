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
func Parser(dbSchemas config.DatabaseSchemas) (model.DBSchemas, error) {
	schema := make(model.DBSchemas)
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
		_, ok := schema[dbSchema.DbAlias]
		if !ok {
			schema[dbSchema.DbAlias] = model.CollectionSchemas{dbSchema.Table: value}
		} else {
			schema[dbSchema.DbAlias][dbSchema.Table] = value
		}
	}
	return schema, nil
}

// GetConstraintName generates constraint name for joint fields
func GetConstraintName(tableName, columnName string) string {
	return fmt.Sprintf("c_%s_%s", tableName, columnName)
}
