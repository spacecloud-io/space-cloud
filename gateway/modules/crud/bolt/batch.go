package bolt

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (b *Bolt) Batch(ctx context.Context, project string, req *model.BatchRequest) ([]int64, error) {
	return nil, nil
}
