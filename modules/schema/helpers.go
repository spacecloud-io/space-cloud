package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetSQLType return sql type
func getSQLType(dbType, typename string) (string, error) {

	switch typename {
	case typeID:
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

	if realFieldStruct.IsPrimary && realFieldStruct.Kind != typeID {
		return errors.New("primary key should be of type ID")
	}

	return nil
}

func (c *creationModule) modifyColumnType() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " MODIFY " + c.FieldKey + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ALTER COLUMN " + c.FieldKey + " TYPE " + c.columnType + " USING (" + c.FieldKey + "::" + c.columnType + ")"
	}
	return ""
}

func (c *creationModule) addNotNull() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " MODIFY " + c.FieldKey + " " + c.columnType + " NOT NULL"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ALTER COLUMN " + c.FieldKey + " SET NOT NULL "
	}
	return ""
}

func (c *creationModule) removeNotNull() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " MODIFY " + c.FieldKey + " " + c.columnType + " NULL"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ALTER COLUMN " + c.FieldKey + " DROP NOT NULL"
	}
	return ""
}

func (c *creationModule) addNewColumn() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD " + c.FieldKey + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD COLUMN " + c.FieldKey + " " + c.columnType
	}
	return ""
}

func (c *creationModule) removeColumn() string {
	return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP COLUMN " + c.FieldKey + ""
}

func (c *creationModule) addPrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD PRIMARY KEY (" + c.FieldKey + ")"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " PRIMARY KEY (" + c.FieldKey + ")"
	}
	return ""
}

func (c *creationModule) removePrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP PRIMARY KEY"
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey
	}
	return ""

}

func (c *creationModule) addUniqueKey() string {
	return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " UNIQUE (" + c.FieldKey + ")"
}

func (c *creationModule) removeUniqueKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP INDEX c_" + c.ColName + "_" + c.FieldKey
	case utils.Postgres:
		return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey
	}
	return ""
}

func (c *creationModule) addForeignKey() string {
	return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " FOREIGN KEY (" + c.FieldKey + ") REFERENCES " + getTableName(c.project, c.realFieldStruct.JointTable.Table, c.removeProjectScope) + "(" + c.realFieldStruct.JointTable.To + ")"
}

func (c *creationModule) removeForeignKey() []string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return []string{"ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP FOREIGN KEY c_" + c.ColName + "_" + c.FieldKey, "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP INDEX c_" + c.ColName + "_" + c.FieldKey}
	case utils.Postgres:
		return []string{"ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey}
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

func (c *creationModule) addField(ctx context.Context) ([]string, error) {
	var queries []string

	if c.columnType != "" {
		// add a new column with data type as columntype
		queries = append(queries, c.addNewColumn())
	}

	if c.realFieldStruct.IsFieldTypeRequired {
		// make the new column not null
		queries = append(queries, c.addNotNull())
	}

	if c.realFieldStruct.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realFieldStruct.IsUnique {
		queries = append(queries, c.addUniqueKey())
	}

	if c.realFieldStruct.IsForeign {
		queries = append(queries, c.addForeignKey())
	}
	return queries, nil
}

func (c *creationModule) removeField() string {
	return c.removeColumn()
}

func (c *creationModule) modifyField(ctx context.Context) ([]string, error) {
	var queries []string

	if !c.realFieldStruct.IsPrimary && c.currentFieldStruct.IsPrimary {
		queries = append(queries, c.removePrimaryKey())
	}

	if !c.realFieldStruct.IsUnique && c.currentFieldStruct.IsUnique {
		queries = append(queries, c.removeUniqueKey())
	}

	if !c.realFieldStruct.IsForeign && c.currentFieldStruct.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
	}

	if c.realFieldStruct.Kind != c.currentFieldStruct.Kind {
		if c.columnType != "" {
			queries = append(queries, c.modifyColumnType())
		}
	}

	if c.realFieldStruct.IsFieldTypeRequired != c.currentFieldStruct.IsFieldTypeRequired {
		if c.realFieldStruct.IsFieldTypeRequired {
			queries = append(queries, c.addNotNull())
		} else {
			queries = append(queries, c.removeNotNull())
		}
	}

	if c.realFieldStruct.IsPrimary && !c.currentFieldStruct.IsPrimary {
		queries = append(queries, c.addPrimaryKey())
	}

	if c.realFieldStruct.IsUnique && !c.currentFieldStruct.IsUnique {
		queries = append(queries, c.addUniqueKey())
	}

	if c.realFieldStruct.IsForeign && !c.currentFieldStruct.IsForeign {
		queries = append(queries, c.addForeignKey())
	}
	return queries, nil
}
