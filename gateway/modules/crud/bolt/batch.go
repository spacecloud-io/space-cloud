package bolt

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Batch performs the provided operations in a single Batch
func (b *Bolt) Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error) {
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Batch operation not supported for selected database", nil, nil)
}
