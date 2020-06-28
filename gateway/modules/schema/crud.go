package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type fieldsToPostProcess struct {
	kind string
	name string
}

// CrudPostProcess unmarshal's the json field in read request
func (s *Schema) CrudPostProcess(ctx context.Context, dbAlias, col string, result interface{}) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	colInfo, ok := s.SchemaDoc[dbAlias]
	if !ok {
		logrus.Errorf("error crud post process in schema module cannot find dbAlias (%s) in schemaDoc", dbAlias)
		return fmt.Errorf("dbAlias (%s) not found in schema doc", dbAlias)
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

	dbType, _ := s.crud.GetDBType(dbAlias)
	var fieldsToProcess []fieldsToPostProcess
	for columnName, columnValue := range tableInfo {
		if columnValue.Kind == model.TypeJSON || columnValue.Kind == model.TypeDateTime || (dbType == string(utils.MySQL) && columnValue.Kind == model.TypeBoolean) {
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
					data, ok := column.([]byte)
					if !ok {
						if column == nil {
							continue
						}
						logrus.Errorf("error crud post process in schema module unable to type assert interface to []byte it is of type (%T) for column (%s)", field.name, doc[field.name])
						return fmt.Errorf("unable to type assert interface to []byte for column (%s)", field.name)
					}
					var v interface{}
					if err := json.Unmarshal(data, &v); err != nil {
						logrus.Errorf("error crud post process in schema module unable unmarshal data (%s)", string(data))
						return fmt.Errorf("unable to unmarshal json data for column (%s)", field.name)
					}
					doc[field.name] = v

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
						doc[field.name] = v.UTC().Format(time.RFC3339)
					case primitive.DateTime:
						doc[field.name] = v.Time().UTC().Format(time.RFC3339)
					}
				}
			}
		}
	}

	return nil
}

// AdjustWhereClause adjusts where clause to take care of types
func (s *Schema) AdjustWhereClause(dbAlias string, dbType utils.DBType, col string, find map[string]interface{}) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Currently this is only required for mongo. Return if its any other database
	if dbType != utils.Mongo {
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
				t, err := time.Parse(time.RFC3339, param)
				if err != nil {
					return fmt.Errorf("invalid string format of datetime (%s) provided for field (%s)", param, k)
				}
				find[k] = t

			case map[string]interface{}:
				for operator, paramInterface := range param {

					// Don't do anything if value is already time.Time
					if _, ok := paramInterface.(time.Time); ok {
						break
					}

					// Check if the value is string
					paramString, ok := paramInterface.(string)
					if !ok {
						return fmt.Errorf("invalid format (%s) of datetime (%v) provided for field (%s)", reflect.TypeOf(paramInterface), paramInterface, k)
					}

					// Try parsing it to time.Time
					t, err := time.Parse(time.RFC3339, paramString)
					if err != nil {
						return fmt.Errorf("invalid string format of datetime (%s) provided for field (%s)", param, k)
					}

					// Store the value
					param[operator] = t
				}
			case time.Time:
				break
			default:
				return fmt.Errorf("invalid format (%s) of datetime (%v) provided for field (%s)", reflect.TypeOf(param), param, k)
			}
		}
	}

	return nil
}
