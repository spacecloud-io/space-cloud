package mgo

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// DescribeTable return a structure of sql table
func (m *Mongo) DescribeTable(ctc context.Context, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	return nil, nil, nil, errors.New("schema operation cannot be performed")
}
