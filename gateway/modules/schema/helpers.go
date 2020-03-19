package schema

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// GetSQLType return sql type
func getSQLType(dbType, typename string) (string, error) {

	switch typename {
	case model.TypeID:
		return "varchar(" + model.SQLTypeIDSize + ")", nil
	case model.TypeString:
		if dbType == string(utils.SQLServer) {
			return "varchar(max)", nil
		}
		return "text", nil
	case model.TypeDateTime:
		if dbType == string(utils.MySQL) {
			return "datetime", nil
		}
		return "timestamp", nil
	case model.TypeBoolean:
		if dbType == string(utils.SQLServer) {
			return "bit", nil
		}
		return "boolean", nil
	case model.TypeFloat:
		return "float", nil
	case model.TypeInteger:
		return "bigint", nil
	default:
		return "", fmt.Errorf("%s type not allowed", typename)
	}
}

func checkErrors(realFieldStruct *model.FieldType) error {
	if realFieldStruct.IsList && !realFieldStruct.IsLinked { // array without directive relation not allowed
		return fmt.Errorf("invalid type for field %s - array type without link directive is not supported in sql creation", realFieldStruct.FieldName)
	}
	if realFieldStruct.Kind == model.TypeObject {
		return fmt.Errorf("invalid type for field %s - object type not supported in sql creation", realFieldStruct.FieldName)
	}

	if realFieldStruct.IsPrimary && !realFieldStruct.IsFieldTypeRequired {
		return errors.New("primary key must be required")
	}

	if realFieldStruct.IsPrimary && realFieldStruct.Kind != model.TypeID {
		return errors.New("primary key should be of type ID")
	}

	if (realFieldStruct.IsUnique || realFieldStruct.IsPrimary || realFieldStruct.IsLinked) && realFieldStruct.IsDefault {
		return errors.New("cannot set default directive with other constraints")
	}

	return nil
}

func (c *creationModule) addNotNull() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	c.currentColumnInfo.IsFieldTypeRequired = true // Mark the field as processed
	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " MODIFY " + c.ColumnName + " " + c.columnType + " NOT NULL"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER COLUMN " + c.ColumnName + " SET NOT NULL"
	case utils.SQLServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " NOT NULL"
	}
	return ""
}

func (c *creationModule) removeNotNull() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " MODIFY " + c.ColumnName + " " + c.columnType + " NULL"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER COLUMN " + c.ColumnName + " DROP NOT NULL"
	case utils.SQLServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " NULL" // adding NULL solves a bug that DateTime type is always not nullable even if (!) is not provided
	}
	return ""
}

func (c *creationModule) addNewColumn() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD " + c.ColumnName + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD COLUMN " + c.ColumnName + " " + c.columnType
	case utils.SQLServer:
		if c.columnType == "timestamp" && !c.realColumnInfo.IsFieldTypeRequired {
			return "ALTER TABLE " + c.project + "." + c.TableName + " ADD " + c.ColumnName + " " + c.columnType + " NULL"
		}

		return "ALTER TABLE " + c.project + "." + c.TableName + " ADD " + c.ColumnName + " " + c.columnType
	}
	return ""
}

func (c *creationModule) removeColumn() []string {
	queries := c.removeDirectives()
	return append(queries, "ALTER TABLE "+getTableName(c.project, c.TableName, c.removeProjectScope)+" DROP COLUMN "+c.ColumnName+"")
}

func (c *creationModule) addPrimaryKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	c.currentColumnInfo.IsPrimary = true // Mark the field as processed
	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD PRIMARY KEY (" + c.ColumnName + ")"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " PRIMARY KEY (" + c.ColumnName + ")"
	case utils.SQLServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " PRIMARY KEY CLUSTERED (" + c.ColumnName + ")"
	}
	return ""
}

func (c *creationModule) removePrimaryKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP PRIMARY KEY"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
	case utils.SQLServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
	}
	return ""

}

func (c *creationModule) addForeignKey() string {
	c.currentColumnInfo.IsForeign = true // Mark the field as processed
	if c.realColumnInfo.JointTable.OnCascade == "CASCADE" {
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " FOREIGN KEY (" + c.ColumnName + ") REFERENCES " + getTableName(c.project, c.realColumnInfo.JointTable.Table, c.removeProjectScope) + " (" + c.realColumnInfo.JointTable.To + ") " + "ON DELETE CASCADE"
	}
	return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " FOREIGN KEY (" + c.ColumnName + ") REFERENCES " + getTableName(c.project, c.realColumnInfo.JointTable.Table, c.removeProjectScope) + " (" + c.realColumnInfo.JointTable.To + ")"
}

func (c *creationModule) typeSwitch() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch v := c.realColumnInfo.Default.(type) {
	case string:
		return "'" + fmt.Sprintf("%v", v) + "'"
	case bool:
		if utils.DBType(dbType) == utils.SQLServer {
			if v {
				return fmt.Sprintf("1")
			}
			return fmt.Sprintf("0")
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (c *creationModule) addDefaultKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	c.currentColumnInfo.IsDefault = true // Mark the field as processed
	c.currentColumnInfo.Default = c.realColumnInfo.Default
	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER " + c.ColumnName + " SET DEFAULT " + c.typeSwitch()

	case utils.SQLServer:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.ColumnName + " DEFAULT " + c.typeSwitch() + " FOR " + c.ColumnName

	case utils.Postgres:

		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER COLUMN " + c.ColumnName + " SET DEFAULT " + c.typeSwitch()
	}
	return ""
}

func (c *creationModule) removeDefaultKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}
	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER " + c.ColumnName + " DROP DEFAULT"

	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER COLUMN " + c.ColumnName + " DROP DEFAULT"
	case utils.SQLServer:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.ColumnName
	}
	return ""
}

func (c *creationModule) removeForeignKey() []string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return nil
	}

	c.currentColumnInfo.IsForeign = false
	switch utils.DBType(dbType) {
	case utils.MySQL:
		return []string{"ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP FOREIGN KEY c_" + c.TableName + "_" + c.ColumnName, "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP INDEX c_" + c.TableName + "_" + c.ColumnName}
	case utils.Postgres:
		return []string{"ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName}
	case utils.SQLServer:
		return []string{"ALTER TABLE " + c.project + "." + c.TableName + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName}
	}
	return nil
}

func addNewTable(project, dbType, realColName string, realColValue model.Fields, removeProjectScope bool) (string, error) {

	var query string
	for realFieldKey, realFieldStruct := range realColValue {

		// Ignore linked fields since these are virtual fields
		if realFieldStruct.IsLinked {
			continue
		}

		if err := checkErrors(realFieldStruct); err != nil {
			return "", err
		}
		sqlType, err := getSQLType(dbType, realFieldStruct.Kind)
		if err != nil {
			return "", nil
		}

		query += realFieldKey + " " + sqlType

		if realFieldStruct.IsPrimary {
			primaryKey := "PRIMARY KEY"
			query += " " + primaryKey
		}

		if realFieldStruct.IsFieldTypeRequired {
			query += " NOT NULL"
		}

		query += " ,"
	}

	return `CREATE TABLE ` + getTableName(project, realColName, removeProjectScope) + ` (` + query[0:len(query)-1] + `);`, nil
}

func getTableName(project, table string, removeProjectScope bool) string {
	if removeProjectScope {
		return table
	}

	return project + "." + table
}

func (c *creationModule) addColumn(dbType string) []string {
	var queries []string

	c.currentColumnInfo = &model.FieldType{
		FieldName: c.realColumnInfo.FieldName,
		Kind:      c.realColumnInfo.Kind,
	}

	if c.columnType != "" {
		// add a new column with data type as columntype
		queries = append(queries, c.addNewColumn())
	}

	if c.realColumnInfo.IsFieldTypeRequired {
		// make the new column not null
		if dbType == string(utils.SQLServer) && c.columnType == "timestamp" {
		} else {
			queries = append(queries, c.addNotNull())
		}
	}

	if c.realColumnInfo.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey())
	}

	if c.realColumnInfo.IsDefault {
		queries = append(queries, c.addDefaultKey())
	}

	return queries
}

// func (c *creationModule) removeField() string {
// 	return c.removeColumn()
// }

func (c *creationModule) modifyColumn() []string {
	var queries []string

	if c.realColumnInfo.IsFieldTypeRequired != c.currentColumnInfo.IsFieldTypeRequired {
		if c.realColumnInfo.IsFieldTypeRequired {
			queries = append(queries, c.addNotNull())
		} else {
			queries = append(queries, c.removeNotNull())
		}
	}

	if !c.realColumnInfo.IsPrimary && c.currentColumnInfo.IsPrimary {
		queries = append(queries, c.removePrimaryKey())
	}

	if !c.realColumnInfo.IsForeign && c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
	}

	if !c.realColumnInfo.IsDefault && c.currentColumnInfo.IsDefault {
		queries = append(queries, c.removeDefaultKey())
	}

	if c.realColumnInfo.IsPrimary && !c.currentColumnInfo.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realColumnInfo.IsForeign && !c.currentColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey())
	}
	// else if (c.realColumnInfo.IsForeign && c.currentColumnInfo.IsForeign) && (c.currentColumnInfo.JointTable.OnCascade != NULL) && (c.currentColumnInfo.JointTable.OnCascade != c.realColumnInfo.JointTable.OnCascade) {
	// 	queries = append(queries, c.removeForeignKey()...)
	// 	queries = append(queries, c.addForeignKey())
	// }

	if c.realColumnInfo.IsDefault && !c.currentColumnInfo.IsDefault {
		queries = append(queries, c.addDefaultKey())
	}

	return queries
}

func (c *creationModule) removeDirectives() []string {
	var queries []string

	if c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
		c.currentColumnInfo.IsForeign = false
	}

	if c.currentColumnInfo.IsDefault {
		queries = append(queries, c.removeDefaultKey())
		c.currentColumnInfo.IsDefault = false
	}

	if c.currentColumnInfo.IsPrimary {
		queries = append(queries, c.removePrimaryKey())
		c.currentColumnInfo.IsPrimary = false
	}

	if c.currentColumnInfo.IsIndex {
		if _, p := c.currentIndexMap[c.currentColumnInfo.IndexInfo.Group]; p {
			queries = append(queries, removeIndex(c.dbAlias, c.project, c.TableName, c.currentColumnInfo.IndexInfo.Group, c.removeProjectScope))
			delete(c.currentIndexMap, c.currentColumnInfo.IndexInfo.Group)
		}
	}

	return queries
}

// modifyColumnType drop the column then creates a new column with provided type
func (c *creationModule) modifyColumnType(dbType string) []string {
	queries := []string{}

	// Remove the column
	queries = append(queries, c.removeColumn()...)

	// Add the column back again
	queries = append(queries, c.addColumn(dbType)...)

	return queries
}

func addIndex(dbType, project, tableName, indexName string, isIndexUnique bool, removeProjectScope bool, mapArray []*model.FieldType) string {
	s := " ("
	for _, schemaFieldType := range mapArray {
		s += schemaFieldType.FieldName + " " + schemaFieldType.IndexInfo.Sort + ", "
	}
	s = strings.TrimSuffix(s, ", ")
	p := ""
	if isIndexUnique {
		p = "CREATE UNIQUE INDEX " + "index__" + tableName + "__" + indexName + " ON " + getTableName(project, tableName, removeProjectScope) + s + ")"
	} else {
		p = "CREATE INDEX " + "index__" + tableName + "__" + indexName + " ON " + getTableName(project, tableName, removeProjectScope) + s + ")"
	}
	return p
}

func removeIndex(dbType, project, tableName, indexName string, removeProjectScope bool) string {

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "DROP INDEX " + "index__" + tableName + "__" + indexName + " ON " + project + "." + tableName
	case utils.SQLServer:
		return "DROP INDEX " + "index__" + tableName + "__" + indexName + " ON " + getTableName(project, tableName, removeProjectScope)
	case utils.Postgres:
		indexname := "index__" + tableName + "__" + indexName
		return "DROP INDEX " + getTableName(project, indexname, removeProjectScope)
	}
	return ""
}

type indexStruct struct {
	IsIndexUnique bool
	IndexMap      []*model.FieldType
}

func getIndexMap(tableInfo model.Fields) (map[string]*indexStruct, error) {
	indexMap := make(map[string]*indexStruct)

	// Iterate over each column of table
	for _, columnInfo := range tableInfo {

		// We are only interested in the columns which have an index on them
		if columnInfo.IsIndex {

			// Append the column to te index map. Make sure we create an empty array if no index by the provided name exists
			value, ok := indexMap[columnInfo.IndexInfo.Group]
			if !ok {
				value = &indexStruct{IndexMap: []*model.FieldType{}}
				indexMap[columnInfo.IndexInfo.Group] = value
			}
			value.IndexMap = append(value.IndexMap, columnInfo)

			// Mark the index group as unique if even on column had the unique tag
			if columnInfo.IsUnique {
				indexMap[columnInfo.IndexInfo.Group].IsIndexUnique = true
			}
		}
	}

	for indexName, indexValue := range indexMap {
		var v indexStore = indexValue.IndexMap
		sort.Stable(v)
		indexValue.IndexMap = v
		for i, column := range indexValue.IndexMap {
			if i+1 != column.IndexInfo.Order {
				return nil, fmt.Errorf("invalid order sequence proveded for index (%s)", indexName)
			}
		}
	}
	return indexMap, nil
}
