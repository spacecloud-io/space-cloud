package bolt

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (b *Bolt) DescribeTable(ctc context.Context, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	return nil, nil, nil, nil
}
