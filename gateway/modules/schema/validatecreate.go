package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SchemaValidator function validates the schema which it gets from module
func (s *Schema) SchemaValidator(col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {

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
				return nil, fmt.Errorf("cannot insert value for a linked field %s", fieldKey)
			}
			continue
		}

		if !ok && fieldValue.IsDefault {
			value = fieldValue.Default
			ok = true
		}

		if fieldValue.Kind == model.TypeID && !ok {
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
func (s *Schema) ValidateCreateOperation(dbAlias, col string, req *model.CreateRequest) error {

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

	collection, ok := s.SchemaDoc[dbAlias]
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
			return fmt.Errorf("invalid document provided for collection (%s:%s)", dbAlias, col)
		}
		newDoc, err := s.SchemaValidator(col, collectionFields, doc)
		if err != nil {
			return err
		}

		v[index] = newDoc
	}

	req.Operation = utils.All
	req.Document = v

	return nil
}
func (s *Schema) checkType(col string, value interface{}, fieldValue *model.FieldType) (interface{}, error) {
	switch v := value.(type) {
	case int:
		// TODO: int64
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v)/1000, 0), nil
		case model.TypeInteger:
			return value, nil
		case model.TypeFloat:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Integer", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case string:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			unitTimeInRFC3339, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("invalid datetime format recieved for field %s in collection %s - use RFC3339 fromat", fieldValue.FieldName, col)
			}
			return unitTimeInRFC3339, nil
		case model.TypeID, model.TypeString:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got String", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case float32, float64:
		switch fieldValue.Kind {
		case model.TypeDateTime:
			return time.Unix(int64(v.(float64))/1000, 0), nil
		case model.TypeFloat:
			return value, nil
		case model.TypeInteger:
			return int64(value.(float64)), nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Float", fieldValue.FieldName, col, fieldValue.Kind)
		}
	case bool:
		switch fieldValue.Kind {
		case model.TypeBoolean:
			return value, nil
		default:
			return nil, fmt.Errorf("invalid type received for field %s in collection %s - wanted %s got Bool", fieldValue.FieldName, col, fieldValue.Kind)
		}

	case time.Time, *time.Time:
		return v, nil

	case map[string]interface{}:
		if fieldValue.Kind == model.TypeJSON {
			data, err := json.Marshal(value)
			if err != nil {
				logrus.Errorf("error checking type in schema module unable to marshal data for field having type json")
				return nil, err
			}
			return string(data), nil
		}
		if fieldValue.Kind != model.TypeObject {
			return nil, fmt.Errorf("invalid type received for field %s in collection %s", fieldValue.FieldName, col)
		}

		return s.SchemaValidator(col, fieldValue.NestedObject, v)

	case []interface{}:
		if !fieldValue.IsList {
			return nil, fmt.Errorf("invalid type (array) received for field %s in collection %s", fieldValue.FieldName, col)
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
