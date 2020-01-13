package bolt

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (b *Bolt) Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error) {
	return nil, fmt.Errorf("aggregate operation not supported for selected database")
}
