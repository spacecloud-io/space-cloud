package schema

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type mockCrudSchemaInterface struct {
	mock.Mock
}

func (m *mockCrudSchemaInterface) GetDBType(dbAlias string) (string, error) {
	c := m.Called(dbAlias)
	return c.String(0), nil
}

func (m *mockCrudSchemaInterface) DescribeTable(ctx context.Context, dbAlias, col string) ([]model.InspectorFieldType, []model.ForeignKeysType, []model.IndexType, error) {
	return nil, nil, nil, nil
}

func (m *mockCrudSchemaInterface) RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error {
	return nil
}
