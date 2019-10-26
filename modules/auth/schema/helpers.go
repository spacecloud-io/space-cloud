package schema

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetSQLType return sql type
func getSQLType(dbType, typeName string) (string, error) {
	switch typeName {
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
	if realFieldStruct.IsList && (realFieldStruct.Directive != directiveRelation) { // array without directive relation not allowed
		return errors.New("schema: array type without relation directive not supported in sql creation")
	}
	if realFieldStruct.Kind == typeObject {
		return errors.New("schema: object type not supported in sql creation")
	}
	if realFieldStruct.Directive == directiveRelation && realFieldStruct.Kind != typeJoin {
		return errors.New("schema : directive relation should contain user defined type got " + realFieldStruct.Kind)
	}
	if realFieldStruct.Kind == typeID && realFieldStruct.Directive != directiveId {
		return errors.New("schema : directive id should have type id")
	} else if realFieldStruct.Kind == typeID && !realFieldStruct.IsFieldTypeRequired {
		return errors.New("schema : id type is must be not nullable (!)")
	}

	return nil
}

func (c *creationModule) modifyColumnType() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ALTER COLUMN " + c.FieldKey + " TYPE " + c.columnType + " USING (" + c.FieldKey + "::" + c.columnType + ")"
	}
	return ""
}

func (c *creationModule) addNotNull() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType + " NOT NULL"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ALTER COLUMN " + c.FieldKey + " SET NOT NULL "
	}
	return ""
}

func (c *creationModule) removeNotNull() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType + " NULL"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ALTER COLUMN " + c.FieldKey + " DROP NOT NULL"
	}
	return ""
}

func (c *creationModule) addNewColumn() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD " + c.FieldKey + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD COLUMN " + c.FieldKey + " " + c.columnType
	}
	return ""
}

func (c *creationModule) removeColumn() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " DROP COLUMN " + c.FieldKey + ""
}

func (c *creationModule) addPrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD PRIMARY KEY (" + c.FieldKey + ")"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " PRIMARY KEY (" + c.FieldKey + ")"
	}
	return ""
}

func (c *creationModule) removePrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP PRIMARY KEY"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey
	}
	return ""

}

func (c *creationModule) addUniqueKey() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " UNIQUE (" + c.FieldKey + ")"
}

func (c *creationModule) removeUniqueKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP INDEX c_" + c.ColName + "_" + c.FieldKey
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey
	}
	return ""
}

func (c *creationModule) addForeignKey() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.ColName + "_" + c.FieldKey + " FOREIGN KEY (" + c.FieldKey + ") REFERENCES " + c.project + "." + c.realFieldStruct.JointTable.TableName + "(" + c.realFieldStruct.JointTable.TableField + ")"
}

func (c *creationModule) removeForeignKey() []string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return []string{"ALTER TABLE " + c.project + "." + c.ColName + " DROP FOREIGN KEY c_" + c.ColName + "_" + c.FieldKey, "ALTER TABLE " + c.project + "." + c.ColName + " DROP INDEX c_" + c.ColName + "_" + c.FieldKey}
	case utils.Postgres:
		return []string{"ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.ColName + "_" + c.FieldKey}
	}
	return nil
}

func addNewTable(project, dbType, realColName string, realColValue schemaField) (string, error) {
	var query string
	var isID bool
	for realFieldKey, realFieldStruct := range realColValue {
		if err := checkErrors(realFieldStruct); err != nil {
			return "", err
		}
		sqlType, err := getSQLType(dbType, realFieldStruct.Kind)
		if err != nil {
			return "", nil
		}
		if realFieldStruct.Kind == typeID && !isID {
			isID = true
			primaryKey := "PRIMARY KEY"
			query += realFieldKey + " " + sqlType + " " + primaryKey + ","
			continue

		}
		query += realFieldKey + " " + sqlType + " ,"
	}

	if !isID {
		return "", errors.New("Schema creation adding new table type id or primary key was not found")
	}
	return `CREATE TABLE ` + project + `.` + realColName + ` (` + query[0:len(query)-1] + `);`, nil
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
	case directiveId:
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
	case directiveId:
		queries = append(queries, c.removePrimaryKey())
	case directiveUnique:
		queries = append(queries, c.removeUniqueKey())
	case directiveRelation:
		queries = append(queries, c.removeForeignKey()...)
	}
	return queries
}