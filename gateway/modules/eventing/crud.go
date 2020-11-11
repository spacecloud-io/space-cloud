package eventing

import (
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) prepareFindObject(req *model.QueueEventRequest) error {
	if req.Type != utils.EventDBUpdate && req.Type != utils.EventDBDelete {
		return nil
	}

	// Get the database event message
	dbRequest := new(model.DatabaseEventMessage)
	if err := mapstructure.Decode(req.Payload, dbRequest); err != nil {
		return err
	}

	// Get the DB type
	dbType, err := m.crud.GetDBType(dbRequest.DBType)
	if err != nil {
		return err
	}

	// Simply return if this is mongo
	if dbType == string(model.Mongo) {
		return nil
	}

	var source map[string]interface{}
	if req.Type == utils.EventDBUpdate {
		source = dbRequest.Doc.(map[string]interface{})
	} else {
		source = dbRequest.Find.(map[string]interface{})
	}

	// Find the primary keys for the table
	primaryKeys := []string{}
	fields, p := m.schema.GetSchema(dbRequest.DBType, dbRequest.Col)
	if p {
		for fieldName, value := range fields {
			if value.IsPrimary {
				primaryKeys = append(primaryKeys, fieldName)
			}
		}
	}

	// Extract primary keys from source and put it in find
	find := map[string]interface{}{}
	for _, key := range primaryKeys {
		if v, p := source[key]; p {
			find[key] = v
		}
	}

	req.Payload.(map[string]interface{})["find"] = find
	return nil
}
