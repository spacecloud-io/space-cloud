package schema

import (
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (s *Schema) schemaValidator(col string, collectionFields SchemaFields, doc map[string]interface{}) (map[string]interface{}, error) {

	mutatedDoc := map[string]interface{}{}
	for fieldKey, fieldValue := range collectionFields {
		// check if key is required
		value, ok := doc[fieldKey]

		if fieldValue.IsLinked {
			if ok {
				return nil, fmt.Errorf("cannot insert value for a linked field %s", fieldKey)
			}
			continue
		}

		if fieldValue.Kind == TypeID && !ok {
			value = ksuid.New().String()
			ok = true
		}

		if fieldValue.IsCreatedAt || fieldValue.IsUpdatedAt {
			mutatedDoc[fieldKey] = time.Now().UTC()
			continue
		}

		if fieldValue.IsFieldTypeRequired {
			if !ok {
				return nil, errors.New("required field " + fieldKey + " from collection " + col + " not present in request")
			}
		}

		// check type
		val, err := s.checkType(col, value, fieldValue)
		if err != nil {
			return nil, err
		}

		mutatedDoc[fieldKey] = val
	}
	return mutatedDoc, nil
}

// ValidateCreateOperation validates schema on create operation
func (s *Schema) ValidateCreateOperation(dbType, col string, req *model.CreateRequest) error {

	if s.SchemaDoc == nil {
		return errors.New("schema not initialized")
	}

	v := make([]interface{}, 0)

	switch t := req.Document.(type) {
	case []interface{}:
		v = t
	case map[string]interface{}:
		v = append(v, t)
	}

	collection, ok := s.SchemaDoc[dbType]
	if !ok {
		return errors.New("No db was found named " + dbType)
	}
	collectionFields, ok := collection[col]
	if !ok {
		return nil
	}

	for index, docTemp := range v {
		doc, ok := docTemp.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid document provided for collection (%s:%s)", dbType, col)
		}
		newDoc, err := s.schemaValidator(col, collectionFields, doc)
		if err != nil {
			return err
		}

		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}
func (s *Schema) checkType(col string, value interface{}, fieldValue *SchemaFieldType) (interface{}, error) {

	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case typeDateTime:
			return time.Unix(int64(v)/1000, 0), nil
		case typeInteger:
			return value, nil
		case typeFloat:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case string:
		switch fieldValue.Kind {
		case typeDateTime:
			unitTimeInRFC3339, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("invalid datetime format recieved for field %s in collection %s - use RFC3339 fromat", fieldValue.FieldName, col)
			}
			return unitTimeInRFC3339, nil
		case TypeID, typeString:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got String", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case float32, float64:
		switch fieldValue.Kind {
		case typeDateTime:
			return time.Unix(int64(v.(float64))/1000, 0), nil
		case typeFloat:
			return value, nil
		case typeInteger:
			return int64(value.(float64)), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Float", fieldValue.FieldName, col, fieldValue.Kind)
		}
	case bool:
		switch fieldValue.Kind {
		case typeBoolean:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Bool", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case map[string]interface{}:
		// TODO: allow this operation for nested insert using links
		if fieldValue.Kind != typeObject {
			return nil, fmt.Errorf("invalid type received for field %s in collection %s", fieldValue.FieldName, col)
		}

		return s.schemaValidator(col, fieldValue.nestedObject, v)

	case []interface{}:
		// TODO: allow this operation for nested insert using links
		if fieldValue.Kind != typeObject {
			return nil, fmt.Errorf("invalid type received for field %s in collection %s", fieldValue.FieldName, col)
		}

		arr := make([]interface{}, len(v))
		for index, value := range v {
			val, err := s.checkType(col, value, fieldValue)
			if err != nil {
				return nil, err
			}
			arr[index] = val
		}
		return arr, nil
	default:
		if !fieldValue.IsFieldTypeRequired {
			return nil, nil
		}

		return nil, fmt.Errorf("no matching type found for field %s in collection %s", fieldValue.FieldName, col)
	}
}
