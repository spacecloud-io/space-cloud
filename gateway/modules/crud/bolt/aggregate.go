package bolt

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (b *Bolt) Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error) {
	return nil, nil
}
