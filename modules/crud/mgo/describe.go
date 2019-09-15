package mgo

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/utils"
)

// DescribeTable return a structure of sql table
func (s *Mongo) DescribeTable(ctx context.Context, project, dbType, col string) ([]utils.FieldType, []utils.ForeignKeysType, error) {
	return nil, nil, errors.New("Scheam operation cannot be performed")
}
