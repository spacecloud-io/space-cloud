package schema

import (
	"errors"
	"time"
)

// ValidateUpdateOperation validates the types of schema during a update request
func (s *Schema) ValidateUpdateOperation(dbType, col string, updateDoc map[string]interface{}) error {
	if len(updateDoc) == 0 {
		return nil
	}
	schemaDb, ok := s.SchemaDoc[dbType]
	if !ok {
		return errors.New(dbType + " not present in Schema")
	}
	SchemaDoc, ok := schemaDb[col]
	if !ok {
		return nil
	}

	for key, doc := range updateDoc {
		switch key {
		case "$set":
			newDoc, err := s.validateSetOperation(doc, SchemaDoc)
			if err != nil {
				return err
			}
			updateDoc[key] = newDoc
		case "$push":
			err := s.validateArrayOperations(doc, SchemaDoc)
			if err != nil {
				return err
			}
		case "$inc", "$min", "$max", "$mul":
			if err := validateMathOperations(doc, SchemaDoc); err != nil {
				return err
			}
		case "$currentDate":
			err := validateDateOperations(doc, SchemaDoc)
			if err != nil {
				return err
			}
		default:
			return errors.New(key + " Update operator not supported")
		}
	}
	return nil
}

func (s *Schema) validateArrayOperations(doc interface{}, SchemaDoc schemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("schema array op wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		SchemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		switch t := fieldValue.(type) {
		case []interface{}:
			if SchemaDocValue.Directive == directiveRelation {
				return errors.New("schema update op array with relation directive not allowed")
			}
			for _, value := range t {
				if _, err := s.checkType(value, SchemaDocValue); err != nil {
					return err
				}
			}
			return nil
		case interface{}:
			if _, err := s.checkType(t, SchemaDocValue); err != nil {
				return err
			}
		default:
			return errors.New("Schema update array op. wrong type ")
		}
	}

	return nil
}

func validateMathOperations(doc interface{}, SchemaDoc schemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("schema math op wrong type passed expecting map[string]interface{}")
	}

	for fieldKey, fieldValue := range v {

		SchemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		switch fieldValue.(type) {
		case int:
			if SchemaDocValue.Kind != typeInteger && SchemaDocValue.Kind != typeFloat {
				return errors.New("Integer : wrong type wanted " + SchemaDocValue.Kind)
			}
			return nil
		case float32, float64:
			if SchemaDocValue.Kind != typeFloat {
				return errors.New("Float : wrong type wanted " + SchemaDocValue.Kind)
			}
			return nil
		default:
			return errors.New("schema update math op. wrong type ")
		}
	}

	return nil
}

func validateDateOperations(doc interface{}, SchemaDoc schemaField) error {

	v, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("schema math op : wrong type passed expecting map[string]interface{}")
	}

	for fieldKey := range v {

		schemaDocValue, ok := SchemaDoc[fieldKey]
		if !ok {
			return errors.New("field not found in schemaField")
		}

		if schemaDocValue.Kind != typeDateTime {
			return errors.New("incorrect data type")
		}
	}

	return nil
}

func (s *Schema) validateSetOperation(doc interface{}, SchemaDoc schemaField) (interface{}, error) {
	v, ok := doc.(map[string]interface{})
	if !ok {
		return nil, errors.New("schema update set op wrong type passed expecting map[string]interface{}")
	}

	newMap := map[string]interface{}{}
	for key, value := range v {
		// check if key present in SchemaDoc
		SchemaDocValue, ok := SchemaDoc[key]
		if !ok {
			return nil, errors.New("schema set op field not found")
		}
		// check type
		newDoc, err := s.checkType(value, SchemaDocValue)
		if err != nil {
			return nil, err
		}
		newMap[key] = newDoc
	}

	for fieldKey, fieldValue := range SchemaDoc {
		if fieldValue.Directive == directiveUpdatedAt {
			newMap[fieldKey] = time.Now().UTC()
		}
	}

	return newMap, nil
}
