package schema

import "github.com/spaceuptech/space-cloud/gateway/model"

type indexStore []*model.FieldType

func (a indexStore) Len() int           { return len(a) }
func (a indexStore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a indexStore) Less(i, j int) bool { return a[i].IndexInfo.Order < a[j].IndexInfo.Order }
