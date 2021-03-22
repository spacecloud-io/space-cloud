package realtime

import "github.com/spaceuptech/space-cloud/gateway/model"

type schemaInterface interface {
	GetSchema(dbAlias, col string) (model.Fields, bool)
}
