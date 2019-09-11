package schema

import (
	"errors"

	"github.com/spaceuptech/space-cloud/utils"
)

// GetSQLType return sql type
func getSQLType(typeName string) (string, error) {
	switch typeName {
	case typeString, typeID:
		return "tinytext", nil
	case typeDateTime:
		return "datetime", nil
	case typeBoolean:
		return "boolean", nil
	case typeFloat:
		return "float", nil
	case typeInteger:
		return "int", nil
	case typeJoin:
		return "", nil
	default:
		return "", errors.New("type not allowed")
	}
}

func (c *creationModule) modifyNotNull() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType + " NOT NULL"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ALTER COLUMN " + c.FieldKey + " SET NOT NULL "
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

func (c *creationModule) addPrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD PRIMARY KEY (" + c.FieldKey + ")"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.FieldKey + " PRIMARY KEY (" + c.FieldKey + ")"
	}
	return ""
}

func (c *creationModule) addUniqueKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD UNIQUE (" + c.FieldKey + ")"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.FieldKey + " UNIQUE (" + c.FieldKey + ")"
	}
	return ""
}

func (c *creationModule) addForeignKey() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " ADD CONSTRAINT c_" + c.FieldKey + " FOREIGN KEY (" + c.FieldKey + ") REFERENCES " + c.project + "." + c.realFieldStruct.JointTable.TableName + "(" + c.realFieldStruct.JointTable.TableField + ")"
}

func (c *creationModule) modifyColumnType() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ALTER COLUMN " + c.FieldKey + " TYPE " + c.columnType
	}
	return ""
}

func (c *creationModule) removeColumn() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " DROP COLUMN " + c.FieldKey + ""
}

func (c *creationModule) removePrimaryKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP PRIMARY KEY"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.FieldKey
	}
	return ""

}
func (c *creationModule) removeUniqueKey() string {
	switch utils.DBType(c.dbType) {
	case utils.MySQL:
		return "ALTER TABLE " + c.project + "." + c.ColName + " ADD UNIQUE (" + c.FieldKey + ")"
	case utils.Postgres:
		return "ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.FieldKey
	}
	return ""
}
func (c *creationModule) removeForeignKey() string {
	return "ALTER TABLE " + c.project + "." + c.ColName + " DROP CONSTRAINT c_" + c.FieldKey
}

func checkErrors(realFieldStruct *schemaFieldType) error {
	if realFieldStruct.IsList {
		return errors.New("Graphql : array type not supported in sql creation")
	}
	if realFieldStruct.Kind == typeObject {
		return errors.New("Graphql : object type not supported in sql creation")
	}
	if realFieldStruct.Directive == directiveRelation && realFieldStruct.Kind != typeJoin {
		return errors.New("Graphql : directive relation should contain user defined type got " + realFieldStruct.Kind)
	}
	if realFieldStruct.Directive == directiveId && realFieldStruct.Kind != typeID {
		return errors.New("Graphql : directive id should have type id")
	}
	return nil
}
