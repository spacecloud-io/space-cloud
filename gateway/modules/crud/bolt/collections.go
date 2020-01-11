package bolt

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (b *Bolt) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	return nil, nil
}

func (b *Bolt) DeleteCollection(ctx context.Context, project, col string) error {
	return nil
}
