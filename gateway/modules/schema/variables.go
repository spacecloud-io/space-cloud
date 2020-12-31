package schema

import "github.com/spaceuptech/space-cloud/gateway/model"

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
