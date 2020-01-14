package bolt

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (b *Bolt) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	return nil, fmt.Errorf("error getting collection operation not supported for selected database")
}

func (b *Bolt) DeleteCollection(ctx context.Context, project, col string) error {
	return fmt.Errorf("error deleting collection operation not supported for selected database")
}
