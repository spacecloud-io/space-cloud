package bolt

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Aggregate performs a bolt db pipeline aggregation
func (b *Bolt) Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error) {
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "aggregate operation not supported for selected database", nil, nil)
}
