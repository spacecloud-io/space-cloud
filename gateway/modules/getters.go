package modules

import "github.com/spaceuptech/space-cloud/gateway/model"

// GetSchemaModule gets the schema module
func (m *Modules) GetSchemaModule() model.SchemaEventingInterface {
	return m.Schema
}
