package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

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
		logrus.Errorf("error crud post process in schema module cannot find collection (%s) in schemaDoc", col)
		return fmt.Errorf("collection (%s) not found in schema doc", col)
	}
	// todo check for array
	docs := make([]interface{}, 0)
	switch v := result.(type) {
	case []interface{}:
		docs = v
	case map[string]interface{}:
		docs = []interface{}{v}
	}
	for columnName, columnValue := range tableInfo {
		if columnValue.Kind == model.TypeJsonb {
			for _, doc := range docs {
				finalDoc := doc.(map[string]interface{})
				data, ok := finalDoc[columnName].([]byte)
				if !ok {
					logrus.Errorf("error crud post process in schema module unable to type assert interface to []byte for column (%s)", columnName)
					return fmt.Errorf("unable to type assert interface to []byte")
				}
				var v interface{}
				if err := json.Unmarshal(data, &v); err != nil {
					logrus.Errorf("error crud post process in schema module unable unmarshal data (%s)", string(data))
					return fmt.Errorf("unable to unmarshal json data")
				}
				finalDoc[columnName] = v
			}
		}
	}
	return nil
}
