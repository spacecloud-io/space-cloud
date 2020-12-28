package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type fieldsToPostProcess struct {
	kind string
	name string
}

// CrudPostProcess unmarshal's the json field in read request
func (s *Schema) CrudPostProcess(ctx context.Context, dbAlias, col string, result interface{}) error {
	if dbAlias != string(model.Mongo) {
		return nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	colInfo, ok := s.SchemaDoc[dbAlias]
	if !ok {
		dbType, ok := s.dbAliasDBTypeMapping[dbAlias]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Schema validation failed, unknown db alias provided (%s)", dbAlias), nil, nil)
		}
		if model.DBType(dbType) == model.Mongo {
			return nil
		}
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unkown db alias (%s) provided to schema module", dbAlias), nil, nil)
	}
	tableInfo, ok := colInfo[col]
	if !ok {
		// Gracefully return if the schema isn't provided
		return nil
	}
	// todo check for array
	docs := make([]interface{}, 0)
	switch v := result.(type) {
	case []interface{}:
		docs = v
	case map[string]interface{}:
		docs = []interface{}{v}
	}

	// dbType, _ := s.crud.GetDBType(dbAlias)
	var fieldsToProcess []fieldsToPostProcess
	for columnName, columnValue := range tableInfo {
		if columnValue.Kind == model.TypeDateTime {
			fieldsToProcess = append(fieldsToProcess, fieldsToPostProcess{kind: columnValue.Kind, name: columnName})
		}
	}

	// Iterate over the docs only if fields need to be post processed
	if len(fieldsToProcess) > 0 {
		for _, temp := range docs {
			doc := temp.(map[string]interface{})

			for _, field := range fieldsToProcess {
				column, ok := doc[field.name]
				if !ok {
					continue
				}

				switch field.kind {
				case model.TypeJSON:
					switch data := column.(type) {
					case []byte:
						var v interface{}
						if err := json.Unmarshal(data, &v); err != nil {
							return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Database contains corrupted json data", err, map[string]interface{}{"type": "[]byte"})
						}
						doc[field.name] = v

					case string:
						var v interface{}
						if err := json.Unmarshal([]byte(data), &v); err != nil {
							return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Database contains corrupted json data", err, map[string]interface{}{"type": "string"})
						}
						doc[field.name] = v
					}

				case model.TypeBoolean:
					switch v := column.(type) {
					case int64:
						if v == int64(1) {
							doc[field.name] = true
						} else {
							doc[field.name] = false
						}
					}

				case model.TypeDateTime:
					switch v := column.(type) {
					case time.Time:
						doc[field.name] = v.UTC().Format(time.RFC3339Nano)
					case primitive.DateTime:
						doc[field.name] = v.Time().UTC().Format(time.RFC3339Nano)
					}
				}
			}
		}
	}

	return nil
}

// AdjustWhereClause adjusts where clause to take care of types
func (s *Schema) AdjustWhereClause(ctx context.Context, dbAlias string, dbType model.DBType, col string, find map[string]interface{}) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Currently this is only required for mongo. Return if its any other database
	if dbType != model.Mongo {
		return nil
	}

	colInfo, ok := s.SchemaDoc[dbAlias]
	if !ok {
		// Gracefully return if the schema isn't provided
		return nil
	}

	tableInfo, ok := colInfo[col]
	if !ok {
		// Gracefully return if the schema isn't provided
		return nil
	}

	for k, v := range find {
		field, p := tableInfo[k]
		if !p {
			continue
		}

		switch field.Kind {
		case model.TypeDateTime:
			switch param := v.(type) {
			case string:
				t, err := time.Parse(time.RFC3339Nano, param)
				if err != nil {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid string format of datetime (%s) provided for field (%s)", param, k), err, nil)
				}
				find[k] = primitive.NewDateTimeFromTime(t)

			case map[string]interface{}:
				for operator, paramInterface := range param {

					// Don't do anything if value is already time.Time
					if t, ok := paramInterface.(time.Time); ok {
						param[operator] = primitive.NewDateTimeFromTime(t)
						continue
					}

					if _, ok := paramInterface.(primitive.DateTime); ok {
						continue
					}

					// Check if the value is string
					paramString, ok := paramInterface.(string)
					if !ok {
						return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid format (%s) of datetime (%v) provided for field (%s)", reflect.TypeOf(paramInterface), paramInterface, k), nil, nil)
					}

					// Try parsing it to time.Time
					t, err := time.Parse(time.RFC3339Nano, paramString)
					if err != nil {
						return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid string format of datetime (%s) provided for field (%s)", param, k), nil, nil)
					}

					// Store the value
					param[operator] = primitive.NewDateTimeFromTime(t)
				}
			case time.Time:
				break
			default:
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid format (%s) of datetime (%v) provided for field (%s)", reflect.TypeOf(param), param, k), nil, nil)
			}
		}
	}

	return nil
}
