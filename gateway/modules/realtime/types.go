package realtime

type schemaInterface interface {
	CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool)
}
