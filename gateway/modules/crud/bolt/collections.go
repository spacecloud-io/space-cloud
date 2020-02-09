package bolt

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetCollections returns collection / tables name of specified database
func (b *Bolt) GetCollections(ctx context.Context, project string) ([]utils.DatabaseCollections, error) {
	return nil, fmt.Errorf("error getting collection operation not supported for selected database")
}

// DeleteCollection deletes collection / tables name of specified database
func (b *Bolt) DeleteCollection(ctx context.Context, project, col string) error {
	return fmt.Errorf("error deleting collection operation not supported for selected database")
}
