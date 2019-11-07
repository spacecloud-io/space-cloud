package schema

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetSQLType return sql type
func getSQLType(dbType, typename string) (string, error) {

	switch typename {
	case typeID, typeJoin:
		return "varchar(50)", nil
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
		return "", errors.New("type not allowed")
	}
}

func checkErrors(realFieldStruct *schemaFieldType) error {
	switch realFieldStruct.Directive {
	case "", directivePrimary, directiveRelation, directiveUnique, directiveCreatedAt, directiveUpdatedAt:
		break
	default:
		//TODO: uncomment after removing id form events_log
		//return errors.New("unknown directive " + realFieldStruct.Directive)
	}
	if realFieldStruct.IsList && (realFieldStruct.Directive != directiveRelation) { // array without directive relation not allowed
		return errors.New("schema: array type without relation directive not supported in sql creation")
	}
	if realFieldStruct.Kind == typeObject {
		return errors.New("schema: object type not supported in sql creation")
	}
	if realFieldStruct.Directive == directiveRelation && realFieldStruct.Kind != typeJoin {
		return errors.New("schema : directive relation should contain user defined type got " + realFieldStruct.Kind)
	}
	if realFieldStruct.Directive == directivePrimary && !realFieldStruct.IsFieldTypeRequired {
		return errors.New("schema directive primary cannot be null require(!)")
	} else if realFieldStruct.Directive == directivePrimary && realFieldStruct.Kind != typeID {
		return errors.New("schema directive primary should have type id")
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
	return "ALTER TABLE " + getTableName(c.project, c.ColName, c.removeProjectScope) + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " FOREIGN KEY (" + c.FieldKey + ") REFERENCES " + getTableName(c.project, c.realFieldStruct.JointTable.TableName, c.removeProjectScope) + "(" + c.realFieldStruct.JointTable.TableField + ")"
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

func addNewTable(project, dbType, realColName string, realColValue schemaField, removeProjectScope bool) (string, error) {
	var query string
	for realFieldKey, realFieldStruct := range realColValue {
		if err := checkErrors(realFieldStruct); err != nil {
			return "", err
		}
		sqlType, err := getSQLType(dbType, realFieldStruct.Kind)
		if err != nil {
			return "", nil
		}

		query += realFieldKey + " " + sqlType

		if realFieldStruct.Directive == directivePrimary {
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
	tempQuery, err := c.addDirective(ctx)
	if err != nil {
		return nil, err
	}
	queries = append(queries, tempQuery...)
	return queries, nil
}

func (c *creationModule) removeField() string {
	return c.removeColumn()
}

func (c *creationModule) modifyField(ctx context.Context) ([]string, error) {
	var queries []string

	if c.realFieldStruct.Directive != c.currentFieldStruct.Directive {
		if c.realFieldStruct.Directive == "" {
			queries = append(queries, c.removeDirective()...)
		}
	}

	if c.realFieldStruct.Kind == typeJoin {
		c.realFieldStruct.Kind = c.realFieldStruct.JointTable.TableName
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
	if c.realFieldStruct.Directive != c.currentFieldStruct.Directive {
		if c.realFieldStruct.Directive != "" {
			tempQuery, err := c.addDirective(ctx)
			if err != nil {
				return nil, err
			}
			queries = append(queries, tempQuery...)
		}
	}
	return queries, nil
}

func (c *creationModule) addDirective(ctx context.Context) ([]string, error) {
	queries := []string{}
	switch c.realFieldStruct.Directive {
	case directivePrimary:
		queries = append(queries, c.addPrimaryKey())
	case directiveUnique:
		queries = append(queries, c.addUniqueKey())
	case directiveRelation:
		queries = append(queries, c.addForeignKey())
	}
	return queries, nil
}

func (c *creationModule) removeDirective() []string {
	queries := []string{}
	switch c.currentFieldStruct.Directive {
	case directivePrimary:
		queries = append(queries, c.removePrimaryKey())
	case directiveUnique:
		queries = append(queries, c.removeUniqueKey())
	case directiveRelation:
		queries = append(queries, c.removeForeignKey()...)
	}
	return queries
}
