package schema

import (
	"fmt"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// ValidateUpdateOperation validates the types of schema during a update request
func (s *Schema) ValidateUpdateOperation(dbAlias, col, op string, updateDoc, find map[string]interface{}) error {
	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := s.SchemaDoc[dbAlias]
	if !ok {
		return fmt.Errorf("%s is not present in schema", dbAlias)
	}
	SchemaDoc, ok := schemaDb[col]
	if !ok {
		return nil
	}

	for key, doc := range updateDoc {
		switch key {
		case "$set":
			newDoc, err := s.validateSetOperation(col, doc, SchemaDoc)
			if err != nil {
				return err
			}
			updateDoc[key] = newDoc
		case "$push":
			err := s.validateArrayOperations(col, doc, SchemaDoc)
			if err != nil {
				return err
			}
		case "$inc", "$min", "$max", "$mul":
			if err := validateMathOperations(col, doc, SchemaDoc); err != nil {
				return err
			}
		case "$currentDate":
			err := validateDateOperations(col, doc, SchemaDoc)
			if err != nil {
				return err
			}
		default:
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

func (s *Schema) validateArrayOperations(col string, doc interface{}, SchemaDoc SchemaFields) error {

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

func validateMathOperations(col string, doc interface{}, SchemaDoc SchemaFields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return fmt.Errorf("field %s from collection %s is not defined in the schema", fieldKey, col)
		}

		switch fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != typeInteger && schemaDocValue.Kind != typeFloat {
				return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldKey, col, schemaDocValue.Kind)
			}
			return nil
		case float32, float64:
			if schemaDocValue.Kind != typeFloat {
				return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Float", fieldKey, col, schemaDocValue.Kind)
			}
			return nil
		default:
			return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind)
		}
	}

	return nil
}

func validateDateOperations(col string, doc interface{}, SchemaDoc SchemaFields) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("document not of type object in collection %s", col)
	}

	for fieldKey := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return fmt.Errorf("field %s from collection %s is not defined in the schema", fieldKey, col)
		}

		if schemaDocValue.Kind != typeDateTime {
			return fmt.Errorf("invalid type received for field %s in collection %s - wanted %s", fieldKey, col, schemaDocValue.Kind)
		}
	}

	return nil
}

func (s *Schema) validateSetOperation(col string, doc interface{}, SchemaDoc SchemaFields) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("document not of type object in collection %s", col)
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// check if key present in SchemaDoc
		SchemaDocValue, ok := SchemaDoc[key]
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
