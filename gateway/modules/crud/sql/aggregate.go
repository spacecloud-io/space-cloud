package sql

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Aggregate performs a mongo db pipeline aggregation
func (s *SQL) Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error) {
	return nil, errors.New("aggregation is not supported for sql databases")
}
