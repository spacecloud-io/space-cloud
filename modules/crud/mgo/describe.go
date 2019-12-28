package mgo

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/utils"
)

// DescribeTable return a structure of sql table
func (m *Mongo) DescribeTable(ctx context.Context, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	return nil, nil, nil, errors.New("schema operation cannot be performed")
}
