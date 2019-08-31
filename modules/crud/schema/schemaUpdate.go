package schema

import (
	"errors"
	"fmt"
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
		return errors.New(col + " Col Not present in Schema of " + dbType)
	}

	var temp interface{}
	var err error
	for key, doc := range updateDoc {
		switch key {
		case "set":
			temp, err = s.validateSetOperation(doc, schemaDoc)
			if err != nil {
				return err
			}
		case "push":
			temp, err = s.validateArrayOperations(doc, schemaDoc)
			if err != nil {
				return err
			}
		case "inc", "min", "max", "mul":
			if err := validateMathOperations(doc, schemaDoc); err != nil {
				return err
			}
		case "currentDate", "currentTimestamp":
			temp, err = validateDateOperations(doc, schemaDoc)
			if err != nil {
				return err
			}
		default:
			return errors.New(key + " Update operator not supported")
		}
		updateDoc[key] = temp
	}
	return nil
}

func (s *Schema) validateArrayOperations(doc interface{}, schemaDoc SchemaField) (interface{}, error) {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, errors.New("Schema math op wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := schemaDoc[fieldKey]
		if !ok {
			return nil, errors.New("field not found in schemaField")
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			arr := make([]interface{}, len(t))
			for index, value := range t {
				newVal, err := s.checkType(value, schemaDocValue)
				if err != nil {
					return nil, err
				}
				arr[index] = newVal
			}
			return arr, nil
		default:
			return nil, errors.New("Schema Update Math Op. Wrong type ")
		}
	}

	return nil, nil
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
			switch schemaDocValue.Kind {
			case TypeInteger:
			default:
				return errors.New("Integer wrong type wanted Integer")
			}
		case float32, float64:
			switch schemaDocValue.Kind {
			case TypeFloat:
			default:
				return errors.New("Float wrong type wanted Float")
			}
		default:
			return errors.New("Schema Update Math Op. Wrong type ")
		}
	}

	return nil
}

func validateDateOperations(doc interface{}, schemaDoc SchemaField) (interface{}, error) {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, errors.New("Schema math op wrong type passed expecting map[string]interface{}")
	}
	newMap := map[string]interface{}{}
	var temp interface{}
	for fieldKey, fieldValue := range v {

		schemaDocValue, ok := schemaDoc[fieldKey]
		if !ok {
			return nil, errors.New("field not found in schemaField")
		}

		temp = fieldValue
		switch t := fieldValue.(type) {
		case int:
			switch schemaDocValue.Kind {
			case TypeDateTime:
				temp = time.Unix(int64(t), 0)
			default:
				return nil, errors.New("Integer wrong type wanted Datetime ")
			}
		case string:
			switch schemaDocValue.Kind {
			case TypeDateTime:
				unitTimeInRFC3339, err := time.Parse(time.RFC3339, t)
				if err != nil {
					return nil, errors.New("String Wrong Date-Time Format")
				}
				temp = unitTimeInRFC3339
			default:
				return nil, errors.New("String wrong type wanted Datetime")
			}
		default:
			return nil, errors.New("Schema Update Math Op. Wrong type ")
		}
		newMap[fieldKey] = temp
	}

	return newMap, nil
}

func (s *Schema) validateSetOperation(doc interface{}, schemaDoc SchemaField) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, errors.New("Schema math op wrong type passed expecting map[string]interface{}")
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// check if key present in schemaDoc if not insert the field
		schemaDocValue, ok := schemaDoc[key]
		if !ok {
			return nil, errors.New("field not found in schemaField")
		}
		fmt.Println("Doc val ", value)
		// check type
		newDoc, err := s.checkType(value, schemaDocValue)
		if err != nil {
			return nil, err
		}
		newMap[key] = newDoc
	}
	return newMap, nil
}
