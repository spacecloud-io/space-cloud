package schema

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// ValidateUpdateOperation validates the types of schema during a update request
func (s *Schema) ValidateUpdateOperation(dbAlias, col, op string, updateDoc, find map[string]interface{}) error {
	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := s.SchemaDoc[dbAlias]
	if !ok {
		logrus.Errorf("error validating update operation in schema module dbAlias (%s) not found in schemaDoc of schema module", dbAlias)
		return fmt.Errorf("%s is not present in schema", dbAlias)
	}
	SchemaDoc, ok := schemaDb[col]
	if !ok {
		logrus.Infof("validating update operation in schema module collection (%s) not found in schemaDoc where dbAlias (%s)", col, dbAlias)
		return nil
	}

	for key, doc := range updateDoc {
		switch key {
		case "$unset":
			return s.validateUnsetOperation(dbAlias, col, doc, SchemaDoc)
		case "$set":
			newDoc, err := s.validateSetOperation(col, doc, SchemaDoc)
			if err != nil {
				logrus.Errorf("error validating set operation in schema module unable to validate (%s) data", key)
				return err
			}
			updateDoc[key] = newDoc
		case "$push":
			err := s.validateArrayOperations(col, doc, SchemaDoc)
			if err != nil {
				logrus.Errorf("error validating array operation in schema module unable to validate (%s) data", key)
				return err
			}
		case "$inc", "$min", "$max", "$mul":
			if err := validateMathOperations(col, doc, SchemaDoc); err != nil {
				logrus.Errorf("error validating math operation in schema module unable to validate (%s) data", key)
				return err
			}
		case "$currentDate":
			err := validateDateOperations(col, doc, SchemaDoc)
			if err != nil {
				logrus.Errorf("error validating date operation in schema module unable to validate (%s) data", key)
				return err
			}
		default:
			logrus.Errorf("error validating update operation in schema module unknown update operator provided (%s)", key)
			return fmt.Errorf("%s update operator is not supported", key)
		}
	}

	// Fill in absent ids and default values
	for fieldName, fieldStruct := range SchemaDoc {
		if op == utils.Upsert && fieldStruct.IsFieldTypeRequired {
			if _, isFieldPresentInFind := find[fieldName]; isFieldPresentInFind || isFieldPresentInUpdate(fieldName, updateDoc) {
				continue
			}
			return fmt.Errorf("required field (%s) not present during upsert", fieldName)
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

func (s *Schema) validateArrayOperations(col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return fmt.Errorf("field %s from collection %s is not defined in the schema", fieldKey, col)
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			if schemaDocValue.IsForeign && !schemaDocValue.IsList {
				return fmt.Errorf("invalid type provided for field %s in collection %s", fieldKey, col)
			}
			for _, value := range t {
				if _, err := s.checkType(col, value, schemaDocValue); err != nil {
					return err
				}
			}
			return nil
		case interface{}:
			if _, err := s.checkType(col, t, schemaDocValue); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid type provided for field %s in collection %s", fieldKey, col)
		}
	}

	return nil
}

func validateMathOperations(col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	for fieldKey, fieldValue := range v {
		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return fmt.Errorf("field %s from collection %s is not defined in the schema", fieldKey, col)
		}
		if schemaDocValue.Kind == model.TypeInteger && reflect.TypeOf(fieldValue).Kind() == reflect.Float64 {
			fieldValue = int(fieldValue.(float64))
		}
		switch fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != model.TypeInteger && schemaDocValue.Kind != model.TypeFloat {
				return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldKey, col, schemaDocValue.Kind)
			}
			return nil
		case float32, float64:
			if schemaDocValue.Kind != model.TypeFloat {
				return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Float", fieldKey, col, schemaDocValue.Kind)
			}
			return nil
		default:
			return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind)
		}
	}

	return nil
}

func validateDateOperations(col string, doc interface{}, SchemaDoc model.Fields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	for fieldKey := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return fmt.Errorf("field %s from collection %s is not defined in the schema", fieldKey, col)
		}

		if schemaDocValue.Kind != model.TypeDateTime {
			return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind)
		}
	}

	return nil
}

func (s *Schema) validateUnsetOperation(dbAlias, col string, doc interface{}, schemaDoc model.Fields) error {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	// Get the db type
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return err
	}

	// For mongo we need to check if the field to be removed is required
	if dbType == string(utils.Mongo) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.IsFieldTypeRequired {
					return fmt.Errorf("cannot use $unset on field which is required")
				}
			}
		}
		return nil
	}

	if dbType == string(utils.Postgres) || dbType == string(utils.MySQL) || dbType == string(utils.SQLServer) {
		for fieldName := range v {
			columnInfo, ok := schemaDoc[strings.Split(fieldName, ".")[0]]
			if ok {
				if columnInfo.Kind == model.TypeJSON {
					return fmt.Errorf("cannot use $unset on field which has type (%s)", model.TypeJSON)
				}
			} else {
				return fmt.Errorf("specified column (%s) doesn't exists in schema of (%s)", fieldName, col)
			}
		}
	}
	return nil
}

func (s *Schema) validateSetOperation(col string, doc interface{}, SchemaDoc model.Fields) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("document not of type object in collection %s", col)
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// We could get a a key with value like `a.b`, where the user intends to set the field `b` inside object `a`. This holds true for working with json
		// types in postgres. However, no such key would be present in the schema. Hence take the top level key to validate the schema
		SchemaDocValue, ok := SchemaDoc[strings.Split(key, ".")[0]]
		if !ok {
			return nil, fmt.Errorf("field %s from collection %s is not defined in the schema", key, col)
		}
		// check type
		newDoc, err := s.checkType(col, value, SchemaDocValue)
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
