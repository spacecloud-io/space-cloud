package schema

import (
	"errors"
	"time"
)

// ValidateUpdateOperation valides the types of schema during a update request
func (s *Schema) ValidateUpdateOperation(dbType, col string, updateDoc map[string]interface{}) error {
	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := s.SchemaDoc[dbType]
	if !ok {
		return errors.New(dbType + " Db Not present in Schema")
	}
	schemaDoc, ok := schemaDb[col]
	if !ok {
		return nil
	}

	for key, doc := range updateDoc {
		switch key {
		case "set":
			newDoc, err := s.validateSetOperation(doc, schemaDoc)
			if err != nil {
				return err
			}
			updateDoc[key] = newDoc
		case "push":
			err := s.validateArrayOperations(doc, schemaDoc)
			if err != nil {
				return err
			}
		case "inc", "min", "max", "mul":
			if err := validateMathOperations(doc, schemaDoc); err != nil {
				return err
			}
		case "currentDate", "currentTimestamp":
			err := validateDateOperations(doc, schemaDoc)
			if err != nil {
				return err
			}
		default:
			return errors.New(key + " Update operator not supported")
		}
	}
	return nil
}

func (s *Schema) validateArrayOperations(doc interface{}, schemaDoc SchemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("Schema math op wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := schemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			for _, value := range t {
				if _, err := s.checkType(value, schemaDocValue); err != nil {
					return err
				}
			}
			return nil
		case interface{}:
			if _, err := s.checkType(t, schemaDocValue); err != nil {
				return err
			}
		default:
			return errors.New("Schema update array op. wrong type ")
		}
	}

	return nil
}

func validateMathOperations(doc interface{}, schemaDoc SchemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("Schema math op wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := schemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		switch fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != TypeInteger {
				return errors.New("Integer : wrong type wanted " + schemaDocValue.Kind)
			}
			return nil
		case float32, float64:
			if schemaDocValue.Kind != TypeFloat {
				return errors.New("Float : wrong type wanted " + schemaDocValue.Kind)
			}
			return nil
		default:
			return errors.New("Schema update math op. wrong type ")
		}
	}

	return nil
}

func validateDateOperations(doc interface{}, schemaDoc SchemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("Schema math op : wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := schemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		switch t := fieldValue.(type) {
		case int:
			if schemaDocValue.Kind != TypeDateTime {
				return errors.New("Integer : wrong type wanted " + schemaDocValue.Kind)
			}
			return nil
		case string:
			if schemaDocValue.Kind != TypeDateTime {
				return errors.New("String : wrong type, wanted " + schemaDocValue.Kind)
			}
			_, err := time.Parse(time.RFC3339, t)
			if err != nil {
				return errors.New("String : wrong date-time format")
			}
			return nil
		default:
			return errors.New("Schema update date op. wrong type ")
		}
	}

	return nil
}

func (s *Schema) validateSetOperation(doc interface{}, schemaDoc SchemaField) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, errors.New("Schema update set op wrong type passed expecting map[string]interface{}")
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// check if key present in schemaDoc if not insert the field
		schemaDocValue, ok := schemaDoc[key]
		if !ok {
			return nil, errors.New("Scheam set op. field not found in schemaField")
		}
		// check type
		newDoc, err := s.checkType(value, schemaDocValue)
		if err != nil {
			return nil, err
		}
		newMap[key] = newDoc
	}
	return newMap, nil
}
