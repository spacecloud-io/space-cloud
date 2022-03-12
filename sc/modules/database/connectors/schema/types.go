package schema

import (
	"context"

	"github.com/spacecloud-io/space-cloud/model"
)

type fieldsToPostProcess struct {
	kind string
	name string
}

type creationModule struct {
	dbType, dbName, TableName, ColumnName, columnType string
	currentIndexMap                                   map[string]*indexStruct
	currentColumnInfo, realColumnInfo                 *model.FieldType
}

type indexStore []*model.TableProperties

func (a indexStore) Len() int           { return len(a) }
func (a indexStore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a indexStore) Less(i, j int) bool { return a[i].Order < a[j].Order }

type primaryKeyStore []*model.FieldType

func (a primaryKeyStore) Len() int      { return len(a) }
func (a primaryKeyStore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a primaryKeyStore) Less(i, j int) bool {
	return a[i].PrimaryKeyInfo.Order < a[j].PrimaryKeyInfo.Order
}

type dbSchemaResponse struct {
	DbAlias   string             `json:"dbAlias"`
	Col       string             `json:"col"`
	Schema    string             `json:"schema,omitempty"`
	SchemaObj model.FieldSchemas `json:"schemaObj,omitempty"`
}

type createSchemaFunc func(ctx context.Context, tableName string, newSchema model.CollectionSchemas) error
