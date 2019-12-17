package schema

import (
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetSQLType return sql type
func getSQLType(dbType, typename string) (string, error) {

	switch typename {
	case TypeID:
		return "varchar(" + sqlTypeIDSize + ")", nil
	case typeString:
		return "text", nil
	case typeDateTime:
		if dbType == string(utils.MySQL) {
			return "datetime", nil
		}
		return "timestamp", nil
	case typeBoolean:
		return "boolean", nil
	case typeFloat:
		return "float", nil
	case typeInteger:
		return "bigint", nil
	default:
		return "", fmt.Errorf("%s type not allowed", typename)
	}
}

func checkErrors(realFieldStruct *SchemaFieldType) error {
	if realFieldStruct.IsList && !realFieldStruct.IsLinked { // array without directive relation not allowed
		return fmt.Errorf("invalid type for field %s - array type without link directive is not supported in sql creation", realFieldStruct.FieldName)
	}
	if realFieldStruct.Kind == typeObject {
		return fmt.Errorf("invalid type for field %s - object type not supported in sql creation", realFieldStruct.FieldName)
	}

	if realFieldStruct.IsPrimary && !realFieldStruct.IsFieldTypeRequired {
		return errors.New("primary key must be required")
	}

	if realFieldStruct.IsPrimary && realFieldStruct.Kind != TypeID {
		return errors.New("primary key should be of type ID")
	}

	return nil
}

func (c *creationModule) addNotNull() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " MODIFY " + c.ColumnName + " " + c.columnType + " NOT NULL"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ALTER COLUMN " + c.ColumnName + " SET NOT NULL"
	case utils.SqlServer:
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
	case utils.SqlServer:
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
	case utils.SqlServer:
		if c.columnType == "timestamp" && !c.realColumnInfo.IsFieldTypeRequired {
			return "ALTER TABLE " + c.project + "." + c.TableName + " ADD " + c.ColumnName + " " + c.columnType + " NULL"
		} else {
			return "ALTER TABLE " + c.project + "." + c.TableName + " ADD " + c.ColumnName + " " + c.columnType
		}
	}
	return ""
}

func (c *creationModule) removeColumn() string {
	return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP COLUMN " + c.ColumnName + ""
}

func (c *creationModule) addPrimaryKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD PRIMARY KEY (" + c.ColumnName + ")"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " PRIMARY KEY (" + c.ColumnName + ")"
	case utils.SqlServer:
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
	case utils.SqlServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
	}
	return ""

}

func (c *creationModule) addUniqueKey() string {
	return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " UNIQUE (" + c.ColumnName + ")"
}

func (c *creationModule) removeUniqueKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP INDEX c_" + c.TableName + "_" + c.ColumnName
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
	case utils.SqlServer:
		return "ALTER TABLE " + c.project + "." + c.TableName + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
	}
	return ""
}

func (c *creationModule) addForeignKey() string {
	return "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " FOREIGN KEY (" + c.ColumnName + ") REFERENCES " + getTableName(c.project, c.realColumnInfo.JointTable.Table, c.removeProjectScope) + " (" + c.realColumnInfo.JointTable.To + ")"
}

func (c *creationModule) removeForeignKey() []string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return nil
	}

	switch utils.DBType(dbType) {
	case utils.MySQL:
		return []string{"ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP FOREIGN KEY c_" + c.TableName + "_" + c.ColumnName, "ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP INDEX c_" + c.TableName + "_" + c.ColumnName}
	case utils.Postgres:
		return []string{"ALTER TABLE " + getTableName(c.project, c.TableName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName}
	case utils.SqlServer:
		return []string{"ALTER TABLE " + c.project + "." + c.TableName + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName}
	}
	return nil
}

func addNewTable(project, dbType, realColName string, realColValue SchemaFields, removeProjectScope bool) (string, error) {

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

func (c *creationModule) addColumn(dbType string) ([]string, error) {
	var queries []string

	if c.columnType != "" {
		// add a new column with data type as columntype
		queries = append(queries, c.addNewColumn())
	}

	if c.realColumnInfo.IsFieldTypeRequired {
		// make the new column not null
		if dbType == string(utils.SqlServer) && c.columnType == "timestamp" {
		} else {
			queries = append(queries, c.addNotNull())
		}
	}

	if c.realColumnInfo.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realColumnInfo.IsUnique {
		queries = append(queries, c.addUniqueKey())
	}

	if c.realColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey())
	}
	return queries, nil
}

func (c *creationModule) modifyColumn() ([]string, error) {
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

	if !c.realColumnInfo.IsUnique && c.currentColumnInfo.IsUnique {
		queries = append(queries, c.removeUniqueKey())
	}

	if !c.realColumnInfo.IsForeign && c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
	}

	if c.realColumnInfo.IsPrimary && !c.currentColumnInfo.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realColumnInfo.IsUnique && !c.currentColumnInfo.IsUnique {
		queries = append(queries, c.addUniqueKey())
	}

	if c.realColumnInfo.IsForeign && !c.currentColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey())
	}

	return queries, nil
}

// modifyColumnType drop the column then creates a new column with provided type
func (c *creationModule) modifyColumnType(dbType string) ([]string, error) {
	queries := []string{}

	if c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
	}
	queries = append(queries, c.removeColumn())

	q, err := c.addColumn(dbType)
	queries = append(queries, q...)
	if err != nil {
		return nil, err
	}

	return queries, nil
}
