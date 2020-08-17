package schema

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/stretchr/testify/mock"
)

type mockCrudSchemaInterface struct {
	mock.Mock
}

func (m *mockCrudSchemaInterface) GetDBType(dbAlias string) (string, error) {
	c := m.Called(dbAlias)
	return c.String(0), nil
}

func (m *mockCrudSchemaInterface) DescribeTable(ctx context.Context, dbAlias, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	return nil, nil, nil, nil
}

func (m *mockCrudSchemaInterface) RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error {
	return nil
}
