package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidateCreateOperation validates req body against provided schema
func ValidateCreateOperation(ctx context.Context, dbAlias, dbType, col string, schemaDoc model.Schemas, req *model.CreateRequest) error {
	if schemaDoc == nil {
		return errors.New("schema not initialized")
	}

	v := make([]interface{}, 0)

	switch t := req.Document.(type) {
	case []interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	collection, ok := schemaDoc[dbAlias]
	if !ok {
		return errors.New("No db was found named " + dbAlias)
	}
	collectionFields, ok := collection[col]
	if !ok {
		return nil
	}

	for index, docTemp := range v {
		doc, ok := docTemp.(map[string]interface{})
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("document provided for collection (%s:%s)", dbAlias, col), nil, nil)
		}
		newDoc, err := Validate(ctx, dbAlias, dbType, col, collectionFields, doc)
		if err != nil {
			return err
		}

		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}

// ValidateUpdateOperation validates the types of schema during a update request
func ValidateUpdateOperation(ctx context.Context, dbAlias, dbType, col, op string, updateDoc, find map[string]interface{}, schemaDoc model.Schemas) error {
	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := schemaDoc[dbAlias]
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to validate update operation in schema module dbAlias (%s) not found in schema module", dbAlias), nil, nil)
	}
	SchemaDoc, ok := schemaDb[col]
	if !ok {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Validating update operation in schema module collection (%s) not found in schemaDoc where dbAlias (%s)", col, dbAlias), nil)
		return nil
	}

	for key, doc := range updateDoc {
		switch key {
		case "$unset":
			return validateUnsetOperation(ctx, dbType, col, doc, SchemaDoc)
		case "$set":
			newDoc, err := validateSetOperation(ctx, dbAlias, dbType, col, doc, SchemaDoc)
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error validating set operation in schema module unable to validate (%s) data", key), err, nil)
			}
			updateDoc[key] = newDoc
		case "$push":
			err := validateArrayOperations(ctx, dbAlias, dbType, col, doc, SchemaDoc)
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error validating array operation in schema module unable to validate (%s) data", key), err, nil)
			}
		case "$inc", "$min", "$max", "$mul":
			if err := validateMathOperations(ctx, col, doc, SchemaDoc); err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error validating math operation in schema module unable to validate (%s) data", key), err, nil)
			}
		case "$currentDate":
			err := validateDateOperations(ctx, col, doc, SchemaDoc)
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error validating date operation in schema module unable to validate (%s) data", key), err, nil)
			}
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to validate update operation unknown update operator (%s) provided", key), nil, nil)
		}
	}

	// Fill in absent ids and default values
	for fieldName, fieldStruct := range SchemaDoc {
		if op == utils.Upsert && fieldStruct.IsFieldTypeRequired {
			if _, isFieldPresentInFind := find[fieldName]; isFieldPresentInFind || isFieldPresentInUpdate(fieldName, updateDoc) {
				continue
			}
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("As per the schema of (%s) field (%s) is mandatory, but it is not present in current upsert operation", col, fieldName), nil, nil)
		}
	}

	return nil
}

// CrudPostProcess unmarshalls the json field in read request
func CrudPostProcess(ctx context.Context, dbAlias, dbType, col string, schemaDoc model.Schemas, result interface{}) error {
	if dbAlias != string(model.Mongo) {
		return nil
	}

	colInfo, ok := schemaDoc[dbAlias]
	if !ok {
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
func AdjustWhereClause(ctx context.Context, dbAlias string, dbType model.DBType, col string, schemaDoc model.Schemas, find map[string]interface{}) error {
	colInfo, ok := schemaDoc[dbAlias]
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
		case model.TypeBoolean:
			if dbType == model.SQLServer {
				switch param := v.(type) {
				case bool:
					if param {
						find[k] = 1
					} else {
						find[k] = 0
					}
				case map[string]interface{}:
					for operator, paramInterface := range param {
						// Check if the value is boolean
						switch t := paramInterface.(type) {
						case []interface{}:
						case bool:
							if t {
								param[operator] = 1
							} else {
								param[operator] = 0
							}
						default:
							return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type (%s) for boolean (%v) provided for field (%s)", reflect.TypeOf(paramInterface), paramInterface, k), nil, nil)
						}
					}
				default:
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type (%s) for boolean (%v) provided for field (%s)", reflect.TypeOf(param), param, k), nil, nil)
				}
			}
		case model.TypeDateTime:
			if dbType == model.Mongo {
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
	}

	return nil
}
