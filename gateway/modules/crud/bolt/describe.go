package bolt

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// DescribeTable return a structure of sql table
func (b *Bolt) DescribeTable(ctx context.Context, col string) ([]model.InspectorFieldType, []model.IndexType, error) {
	return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Describe table operation not supported for selected database", nil, nil)
}
