package bolt

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// DescribeTable return a structure of sql table
func (b *Bolt) DescribeTable(ctc context.Context, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	return nil, nil, nil, fmt.Errorf("error describing table operation not supported for selected database")
}
