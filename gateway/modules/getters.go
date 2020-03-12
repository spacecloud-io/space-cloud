package modules

import "github.com/spaceuptech/space-cloud/gateway/model"

func (m *Modules) GetSchemaModule() model.SchemaEventingInterface {
	return m.Schema
}
