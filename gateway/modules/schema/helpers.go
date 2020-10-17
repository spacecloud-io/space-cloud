package schema

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// GetSQLType return sql type
func getSQLType(ctx context.Context, maxIDSize int, dbType, typename string) (string, error) {

	switch typename {
	case model.TypeID:
		return fmt.Sprintf("varchar(%d)", maxIDSize), nil
	case model.TypeString:
		if dbType == string(model.SQLServer) {
			return "varchar(max)", nil
		}
		return "text", nil
	case model.TypeDateTime:
		switch dbType {
		case string(model.MySQL):
			return "datetime", nil
		case string(model.SQLServer):
			return "datetimeoffset", nil
		default:
			return "timestamp", nil
		}
	case model.TypeBoolean:
		if dbType == string(model.SQLServer) {
			return "bit", nil
		}
		return "boolean", nil
	case model.TypeFloat:
		return "float", nil
	case model.TypeInteger:
		return "bigint", nil
	case model.TypeJSON:
		switch dbType {
		case string(model.Postgres):
			return "jsonb", nil
		case string(model.MySQL):
			return "json", nil
		default:
			return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("json not supported for database %s", dbType), nil, nil)
		}
	default:
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid schema type (%s) provided", typename), fmt.Errorf("%s type not allowed", typename), nil)
	}
}

func checkErrors(ctx context.Context, realFieldStruct *model.FieldType) error {
	if realFieldStruct.IsList && !realFieldStruct.IsLinked { // array without directive relation not allowed
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("invalid type for field %s - array type without link directive is not supported in sql creation", realFieldStruct.FieldName), nil)
	}
	if realFieldStruct.Kind == model.TypeObject {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("invalid type for field %s - object type not supported in sql creation", realFieldStruct.FieldName), nil)
	}

	if realFieldStruct.IsPrimary && !realFieldStruct.IsFieldTypeRequired {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("primary key must be required"), nil)
	}

	if realFieldStruct.IsPrimary && realFieldStruct.Kind != model.TypeID {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("primary key should be of type ID"), nil)
	}

	if realFieldStruct.Kind == model.TypeJSON && (realFieldStruct.IsUnique || realFieldStruct.IsPrimary || realFieldStruct.IsLinked || realFieldStruct.IsIndex) {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("cannot set index with type json"), nil)
	}

	if (realFieldStruct.IsUnique || realFieldStruct.IsPrimary || realFieldStruct.IsLinked) && realFieldStruct.IsDefault {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid schema provided", fmt.Errorf("cannot set default directive with other constraints"), nil)
	}

	return nil
}

func (c *creationModule) addNotNull() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	c.currentColumnInfo.IsFieldTypeRequired = true // Mark the field as processed
	switch model.DBType(dbType) {
	case model.MySQL:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " MODIFY " + c.ColumnName + " " + c.columnType + " NOT NULL"
	case model.Postgres:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " SET NOT NULL"
	case model.SQLServer:
		if strings.HasPrefix(c.columnType, "varchar") {
			return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " collate Latin1_General_CS_AS NOT NULL"
		}
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " NOT NULL"
	}
	return ""
}

func (c *creationModule) removeNotNull() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch model.DBType(dbType) {
	case model.MySQL:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " MODIFY " + c.ColumnName + " " + c.columnType + " NULL"
	case model.Postgres:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " DROP NOT NULL"
	case model.SQLServer:
		if strings.HasPrefix(c.columnType, "varchar") {
			return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " collate Latin1_General_CS_AS NULL"
		}
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " " + c.columnType + " NULL" // adding NULL solves a bug that DateTime type is always not nullable even if (!) is not provided
	}
	return ""
}

func (c *creationModule) addNewColumn() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}

	switch model.DBType(dbType) {
	case model.MySQL:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD " + c.ColumnName + " " + c.columnType
	case model.Postgres:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD COLUMN " + c.ColumnName + " " + c.columnType
	case model.SQLServer:
		if c.columnType == "timestamp" && !c.realColumnInfo.IsFieldTypeRequired {
			return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD " + c.ColumnName + " " + c.columnType + " NULL"
		}
		if strings.HasPrefix(c.columnType, "varchar") {
			return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD " + c.ColumnName + " " + c.columnType + " collate Latin1_General_CS_AS"
		}

		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD " + c.ColumnName + " " + c.columnType
	}
	return ""
}

func (c *creationModule) removeColumn(dbType string) []string {
	queries := c.removeDirectives(dbType)
	return append(queries, "ALTER TABLE "+c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName)+" DROP COLUMN "+c.ColumnName+"")
}

// func (c *creationModule) addPrimaryKey() string {
// 	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
// 	if err != nil {
// 		return ""
// 	}
//
// 	c.currentColumnInfo.IsPrimary = true // Mark the field as processed
// 	switch utils.DBType(dbType) {
// 	case utils.MySQL:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD PRIMARY KEY (" + c.ColumnName + ")"
// 	case utils.Postgres:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " PRIMARY KEY (" + c.ColumnName + ")"
// 	case utils.SQLServer:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD CONSTRAINT c_" + c.TableName + "_" + c.ColumnName + " PRIMARY KEY CLUSTERED (" + c.ColumnName + ")"
// 	}
// 	return ""
// }

// func (c *creationModule) removePrimaryKey() string {
// 	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
// 	if err != nil {
// 		return ""
// 	}
//
// 	switch utils.DBType(dbType) {
// 	case utils.MySQL:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP PRIMARY KEY"
// 	case utils.Postgres:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
// 	case utils.SQLServer:
// 		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP CONSTRAINT c_" + c.TableName + "_" + c.ColumnName
// 	}
// 	return ""
//
// }

func (c *creationModule) addForeignKey(dbType string) string {
	c.currentColumnInfo.IsForeign = true // Mark the field as processed
	if c.currentColumnInfo.JointTable == nil {
		c.currentColumnInfo.JointTable = &model.TableProperties{}
	}
	c.currentColumnInfo.JointTable.OnDelete = c.realColumnInfo.JointTable.OnDelete

	if c.realColumnInfo.JointTable.OnDelete == "CASCADE" {
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD CONSTRAINT " + c.realColumnInfo.JointTable.ConstraintName + " FOREIGN KEY (" + c.ColumnName + ") REFERENCES " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.realColumnInfo.JointTable.Table) + " (" + c.realColumnInfo.JointTable.To + ") " + "ON DELETE CASCADE"
	}
	return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD CONSTRAINT " + c.realColumnInfo.JointTable.ConstraintName + " FOREIGN KEY (" + c.ColumnName + ") REFERENCES " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.realColumnInfo.JointTable.Table) + " (" + c.realColumnInfo.JointTable.To + ")"
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
		if model.DBType(dbType) == model.SQLServer {
			if v {
				return "1"
			}
			return "0"
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
	switch model.DBType(dbType) {
	case model.MySQL:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER " + c.ColumnName + " SET DEFAULT " + c.typeSwitch()
	case model.SQLServer:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ADD CONSTRAINT c_" + c.ColumnName + " DEFAULT " + c.typeSwitch() + " FOR " + c.ColumnName
	case model.Postgres:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " SET DEFAULT " + c.typeSwitch()
	}
	return ""
}

func (c *creationModule) removeDefaultKey() string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return ""
	}
	switch model.DBType(dbType) {
	case model.MySQL:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER " + c.ColumnName + " DROP DEFAULT"
	case model.Postgres:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " ALTER COLUMN " + c.ColumnName + " DROP DEFAULT"
	case model.SQLServer:
		return "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP CONSTRAINT c_" + c.ColumnName
	}
	return ""
}

func (c *creationModule) removeForeignKey() []string {
	dbType, err := c.schemaModule.crud.GetDBType(c.dbAlias)
	if err != nil {
		return nil
	}

	c.currentColumnInfo.IsForeign = false
	switch model.DBType(dbType) {
	case model.MySQL:
		return []string{"ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP FOREIGN KEY " + c.currentColumnInfo.JointTable.ConstraintName, "ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP INDEX " + c.currentColumnInfo.JointTable.ConstraintName}
	case model.Postgres:
		return []string{"ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP CONSTRAINT " + c.currentColumnInfo.JointTable.ConstraintName}
	case model.SQLServer:
		return []string{"ALTER TABLE " + c.schemaModule.getTableName(dbType, c.logicalDBName, c.TableName) + " DROP CONSTRAINT " + c.currentColumnInfo.JointTable.ConstraintName}
	}
	return nil
}

func (s *Schema) addNewTable(ctx context.Context, logicalDBName, dbType, dbAlias, realColName string, realColValue model.Fields) (string, error) {

	var query, primaryKeyQuery string
	doesPrimaryKeyExists := false
	for realFieldKey, realFieldStruct := range realColValue {

		// Ignore linked fields since these are virtual fields
		if realFieldStruct.IsLinked {
			continue
		}
		if err := checkErrors(ctx, realFieldStruct); err != nil {
			return "", err
		}
		sqlType, err := getSQLType(ctx, realFieldStruct.TypeIDSize, dbType, realFieldStruct.Kind)
		if err != nil {
			return "", nil
		}

		if realFieldStruct.IsPrimary {
			doesPrimaryKeyExists = true
			if (model.DBType(dbType) == model.SQLServer) && (strings.HasPrefix(sqlType, "varchar")) {
				primaryKeyQuery = realFieldKey + " " + sqlType + " collate Latin1_General_CS_AS PRIMARY KEY NOT NULL, "
				continue
			}
			primaryKeyQuery = realFieldKey + " " + sqlType + " PRIMARY KEY NOT NULL, "
			continue
		}

		query += realFieldKey + " " + sqlType

		if (model.DBType(dbType) == model.SQLServer) && (strings.HasPrefix(sqlType, "varchar")) {
			query += " collate Latin1_General_CS_AS"
		}

		if realFieldStruct.IsFieldTypeRequired {
			query += " NOT NULL"
		}

		query += " ,"
	}
	if !doesPrimaryKeyExists {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Primary key not found, make sure there is a primary key on a field with type (ID)", nil, nil)
	}
	if model.DBType(dbType) == model.MySQL {
		return `CREATE TABLE ` + s.getTableName(dbType, logicalDBName, realColName) + ` (` + primaryKeyQuery + strings.TrimSuffix(query, " ,") + `) COLLATE Latin1_General_CS;`, nil
	}
	return `CREATE TABLE ` + s.getTableName(dbType, logicalDBName, realColName) + ` (` + primaryKeyQuery + strings.TrimSuffix(query, " ,") + `);`, nil
}

func (s *Schema) getTableName(dbType, logicalDBName, table string) string {
	switch model.DBType(dbType) {
	case model.Postgres, model.SQLServer:
		return fmt.Sprintf("%s.%s", logicalDBName, table)
	}
	return table
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
		if dbType == string(model.SQLServer) && c.columnType == "timestamp" {
		} else {
			queries = append(queries, c.addNotNull())
		}
	}

	// if c.realColumnInfo.IsPrimary {
	// 	queries = append(queries, c.addPrimaryKey())
	// }

	if c.realColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey(dbType))
	}

	if c.realColumnInfo.IsDefault {
		queries = append(queries, c.addDefaultKey())
	}

	return queries
}

func (c *creationModule) modifyColumn(dbType string) []string {
	var queries []string

	if c.realColumnInfo.IsFieldTypeRequired != c.currentColumnInfo.IsFieldTypeRequired {
		if c.realColumnInfo.IsFieldTypeRequired {
			queries = append(queries, c.addNotNull())
		} else {
			queries = append(queries, c.removeNotNull())
		}
	}

	// if !c.realColumnInfo.IsPrimary && c.currentColumnInfo.IsPrimary {
	// 	queries = append(queries, c.removePrimaryKey())
	// }

	if !c.realColumnInfo.IsForeign && c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
	}

	if !c.realColumnInfo.IsDefault && c.currentColumnInfo.IsDefault {
		queries = append(queries, c.removeDefaultKey())
	}

	// if c.realColumnInfo.IsPrimary && !c.currentColumnInfo.IsPrimary {
	// 	queries = append(queries, c.addPrimaryKey())
	// }

	if c.realColumnInfo.IsForeign && !c.currentColumnInfo.IsForeign {
		queries = append(queries, c.addForeignKey(dbType))
	} else if c.realColumnInfo.IsForeign && c.currentColumnInfo.IsForeign && c.currentColumnInfo.JointTable.OnDelete != c.realColumnInfo.JointTable.OnDelete {
		queries = append(queries, c.removeForeignKey()...)
		queries = append(queries, c.addForeignKey(dbType))
	}

	if c.realColumnInfo.IsDefault && !c.currentColumnInfo.IsDefault {
		queries = append(queries, c.addDefaultKey())
	}

	return queries
}

func (c *creationModule) removeDirectives(dbType string) []string {
	var queries []string

	if c.currentColumnInfo.IsForeign {
		queries = append(queries, c.removeForeignKey()...)
		c.currentColumnInfo.IsForeign = false
	}

	if c.currentColumnInfo.IsDefault {
		queries = append(queries, c.removeDefaultKey())
		c.currentColumnInfo.IsDefault = false
	}

	// if c.currentColumnInfo.IsPrimary {
	// 	queries = append(queries, c.removePrimaryKey())
	// 	c.currentColumnInfo.IsPrimary = false
	// }

	if c.currentColumnInfo.IsIndex {
		if _, p := c.currentIndexMap[c.currentColumnInfo.IndexInfo.Group]; p {
			queries = append(queries, c.schemaModule.removeIndex(dbType, c.dbAlias, c.logicalDBName, c.TableName, c.currentColumnInfo.IndexInfo.ConstraintName))
			delete(c.currentIndexMap, c.currentColumnInfo.IndexInfo.Group)
		}
	}

	return queries
}

// modifyColumnType drop the column then creates a new column with provided type
func (c *creationModule) modifyColumnType(dbType string) []string {
	queries := []string{}

	// Remove the column
	queries = append(queries, c.removeColumn(dbType)...)

	// Add the column back again
	queries = append(queries, c.addColumn(dbType)...)

	return queries
}

func (s *Schema) addIndex(dbType, dbAlias, logicalDBName, tableName, indexName string, isIndexUnique bool, mapArray []*model.FieldType) string {
	a := " ("
	for _, schemaFieldType := range mapArray {
		a += schemaFieldType.FieldName + " " + schemaFieldType.IndexInfo.Sort + ", "
	}
	a = strings.TrimSuffix(a, ", ")
	p := ""
	if isIndexUnique {
		p = "CREATE UNIQUE INDEX " + getIndexName(tableName, indexName) + " ON " + s.getTableName(dbType, logicalDBName, tableName) + a + ")"
	} else {
		p = "CREATE INDEX " + getIndexName(tableName, indexName) + " ON " + s.getTableName(dbType, logicalDBName, tableName) + a + ")"
	}
	return p
}

func (s *Schema) removeIndex(dbType, dbAlias, logicalDBName, tableName, indexName string) string {

	switch model.DBType(dbType) {
	case model.MySQL:
		return "DROP INDEX " + indexName + " ON " + s.getTableName(dbType, logicalDBName, tableName)
	case model.SQLServer:
		return "DROP INDEX " + indexName + " ON " + s.getTableName(dbType, logicalDBName, tableName)
	case model.Postgres:
		indexname := indexName
		return "DROP INDEX " + s.getTableName(dbType, logicalDBName, indexname)
	}
	return ""
}

func getIndexName(tableName, indexName string) string {
	return fmt.Sprintf("index__%s__%s", tableName, indexName)
}

func getConstraintName(tableName, columnName string) string {
	return fmt.Sprintf("c_%s_%s", tableName, columnName)
}

type indexStruct struct {
	IsIndexUnique bool
	IndexMap      []*model.FieldType
	IndexName     string
}

func getIndexMap(ctx context.Context, tableInfo model.Fields) (map[string]*indexStruct, error) {
	indexMap := make(map[string]*indexStruct)

	// Iterate over each column of table
	for _, columnInfo := range tableInfo {

		// We are only interested in the columns which have an index on them
		if columnInfo.IsIndex {

			// Append the column to te index map. Make sure we create an empty array if no index by the provided name exists
			value, ok := indexMap[columnInfo.IndexInfo.Group]
			if !ok {
				value = &indexStruct{IndexMap: []*model.FieldType{}, IndexName: columnInfo.IndexInfo.ConstraintName}
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
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid order sequence proveded for index (%s)", indexName), nil, nil)
			}
		}
	}
	return indexMap, nil
}

func (s *Schema) getSchemaResponse(ctx context.Context, format, dbName, tableName string, ignoreForeignCheck bool, alreadyAddedTables map[string]bool, schemaResponse *[]interface{}) error {
	_, ok := alreadyAddedTables[getKeyName(dbName, tableName)]
	if ok {
		return nil
	}

	resourceID := config.GenerateResourceID(s.clusterID, s.project, config.ResourceDatabaseSchema, dbName, tableName)
	dbSchema, ok := s.dbSchemas[resourceID]
	if !ok {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("collection (%s) not present in config for dbAlias (%s) )", dbName, tableName), nil, nil)
	}

	collectionInfo, _ := s.GetSchema(dbName, tableName)
	for _, fieldInfo := range collectionInfo {
		if !ignoreForeignCheck && fieldInfo.IsForeign {
			_, ok := alreadyAddedTables[getKeyName(dbName, tableName)]
			if ok {
				continue
			}
			if err := s.getSchemaResponse(ctx, format, dbName, fieldInfo.JointTable.Table, ignoreForeignCheck, alreadyAddedTables, schemaResponse); err != nil {
				return err
			}
		}
	}
	alreadyAddedTables[getKeyName(dbName, tableName)] = true
	if format == "json" {
		*schemaResponse = append(*schemaResponse, dbSchemaResponse{DbAlias: dbName, Col: tableName, SchemaObj: collectionInfo})
	} else {
		*schemaResponse = append(*schemaResponse, dbSchemaResponse{DbAlias: dbName, Col: tableName, Schema: dbSchema.Schema})
	}
	return nil
}

func getKeyName(dbName, key string) string {
	return fmt.Sprintf("%s-%s", dbName, key)
}

type dbSchemaResponse struct {
	DbAlias   string       `json:"dbAlias"`
	Col       string       `json:"col"`
	Schema    string       `json:"schema,omitempty"`
	SchemaObj model.Fields `json:"schemaObj,omitempty"`
}
