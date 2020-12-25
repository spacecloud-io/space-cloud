package schema

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// ValidateUpdateOperation validates the types of schema during a update request
func (s *Schema) ValidateUpdateOperation(ctx context.Context, dbAlias, col, op string, updateDoc, find map[string]interface{}) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := s.SchemaDoc[dbAlias]
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
			return s.validateUnsetOperation(ctx, dbAlias, col, doc, SchemaDoc)
		case "$set":
			newDoc, err := s.validateSetOperation(ctx, dbAlias, col, doc, SchemaDoc)
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error validating set operation in schema module unable to validate (%s) data", key), err, nil)
			}
			updateDoc[key] = newDoc
		case "$push":
			err := s.validateArrayOperations(ctx, dbAlias, col, doc, SchemaDoc)
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

func isFieldPresentInUpdate(field string, updateDoc map[string]interface{}) bool {
	for _, operatorTemp := range updateDoc {
		operator := operatorTemp.(map[string]interface{})
		if _, p := operator[field]; p {
			return true
		}
	}

	return false
}

func (s *Schema) validateArrayOperations(ctx context.Context, dbAlias, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection %s", col), nil, nil)
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			if schemaDocValue.IsForeign && !schemaDocValue.IsList {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field %s in collection %s", fieldKey, col), nil, nil)
			}
			for _, value := range t {
				if _, err := s.checkType(ctx, dbAlias, col, value, schemaDocValue); err != nil {
					return err
				}
			}
			return nil
		case interface{}:
			if _, err := s.checkType(ctx, dbAlias, col, t, schemaDocValue); err != nil {
				return err
			}
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type provided for field %s in collection %s", fieldKey, col), nil, nil)
		}
	}

	return nil
}

func validateMathOperations(ctx context.Context, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	for fieldKey, fieldValue := range v {
		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}
		if schemaDocValue.Kind == model.TypeInteger && reflect.TypeOf(fieldValue).Kind() == reflect.Float64 {
			fieldValue = int(fieldValue.(float64))
		}
		switch fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != model.TypeInteger && schemaDocValue.Kind != model.TypeFloat {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s got Integer", fieldKey, col, schemaDocValue.Kind), nil, nil)
			}
			return nil
		case float32, float64:
			if schemaDocValue.Kind != model.TypeFloat {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s got Float", fieldKey, col, schemaDocValue.Kind), nil, nil)
			}
			return nil
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind), nil, nil)
		}
	}

	return nil
}

func validateDateOperations(ctx context.Context, col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	for fieldKey := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field %s from collection %s is not defined in the schema", fieldKey, col), nil, nil)
		}

		if schemaDocValue.Kind != model.TypeDateTime {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind), nil, nil)
		}
	}

	return nil
}

func (s *Schema) validateUnsetOperation(ctx context.Context, dbAlias, col string, doc interface{}, schemaDoc model.Fields) error {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	// Get the db type
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return err
	}

	// For mongo we need to check if the field to be removed is required
	if dbType == string(model.Mongo) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.IsFieldTypeRequired {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot use $unset on field which is required/mandatory", nil, nil)
				}
			}
		}
		return nil
	}

	if dbType == string(model.Postgres) || dbType == string(model.MySQL) || dbType == string(model.SQLServer) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.Kind == model.TypeJSON {
					return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot use $unset on field which has type (%s)", model.TypeJSON), nil, nil)
				}
			} else {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field (%s) doesn't exists in schema of (%s)", fieldName, col), nil, nil)
			}
		}
	}
	return nil
}

func (s *Schema) validateSetOperation(ctx context.Context, dbAlias, col string, doc interface{}, SchemaDoc model.Fields) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Document not of type object in collection (%s)", col), nil, nil)
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// We could get a a key with value like `a.b`, where the user intends to set the field `b` inside object `a`. This holds true for working with json
		// types in postgres. However, no such key would be present in the schema. Hence take the top level key to validate the schema
		SchemaDocValue, ok := SchemaDoc[strings.Split(key, ".")[0]]
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field (%s) from collection (%s) is not defined in the schema", key, col), nil, nil)
		}
		// check type
		newDoc, err := s.checkType(ctx, dbAlias, col, value, SchemaDocValue)
		if err != nil {
			return nil, err
		}
		newMap[key] = newDoc
	}

	for fieldKey, fieldValue := range SchemaDoc {
		if fieldValue.IsUpdatedAt {
			newMap[fieldKey] = time.Now().UTC()
		}
	}

	return newMap, nil
}
