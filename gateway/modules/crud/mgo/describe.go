package mgo

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// DescribeTable return a structure of sql table
func (m *Mongo) DescribeTable(ctc context.Context, col string) ([]model.InspectorFieldType, []model.ForeignKeysType, []model.IndexType, error) {
	return nil, nil, nil, errors.New("schema operation cannot be performed")
}
