package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/model"
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

	var fieldsToProcess []fieldsToPostProcess
	for columnName, columnValue := range tableInfo {
		if columnValue.Kind == model.TypeJSON || columnValue.Kind == model.TypeDateTime {
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
						logrus.Errorf("error crud post process in schema module unable to type assert interface to []byte it is of type (%T) for column (%s)", field.name, doc[field.name])
						return fmt.Errorf("unable to type assert interface to []byte for column (%s)", field.name)
					}
					var v interface{}
					if err := json.Unmarshal(data, &v); err != nil {
						logrus.Errorf("error crud post process in schema module unable unmarshal data (%s)", string(data))
						return fmt.Errorf("unable to unmarshal json data for column (%s)", field.name)
					}
					doc[field.name] = v

				case model.TypeDateTime:
					switch v := column.(type) {
					case time.Time:
						doc[field.name] = v.Format(time.RFC3339)
					}
				}
			}
		}
	}

	return nil
}
