package bolt

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Batch performs the provided operations in a single Batch
func (b *Bolt) Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error) {
	return nil, fmt.Errorf("batch operation not supported for selected database")
}
