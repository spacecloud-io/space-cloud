package schema

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type creationModule struct {
	dbAlias, logicalDBName, TableName, ColumnName, columnType string
	currentIndexMap                                           map[string]*indexStruct
	currentColumnInfo, realColumnInfo                         *model.FieldType
	schemaModule                                              *Schema
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbAlias, tableName, logicalDBName string, parsedSchema model.Type) error {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return err
	}

	// Return gracefully if db type is mongo
	if dbType == string(utils.Mongo) || dbType == string(utils.EmbeddedDB) {
		return nil
	}

	currentSchema, err := s.Inspector(ctx, dbAlias, dbType, logicalDBName, tableName)
	if err != nil {
		logrus.Debugln("Schema Inspector Error", err)
	}

	queries, err := s.generateCreationQueries(ctx, dbAlias, tableName, logicalDBName, parsedSchema, currentSchema)
	if err != nil {
		return err
	}
	return s.crud.RawBatch(ctx, dbAlias, queries)
}

func (s *Schema) generateCreationQueries(ctx context.Context, dbAlias, tableName, logicalDBName string, parsedSchema model.Type, currentSchema model.Collection) ([]string, error) {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return nil, err
	}

	// Return nil, if no tables are present in schema
	if len(parsedSchema[dbAlias]) == 0 {
		return nil, nil
	}

	realSchema := parsedSchema[dbAlias]
	batchedQueries := []string{}

	realTableName := tableName
	realTableInfo, p1 := realSchema[realTableName]
	if !p1 {
		if _, p2 := currentSchema[realTableName]; p2 {
			return nil, nil
		}

		return nil, errors.New("Schema not provided for table: " + tableName)
	}

	// check if table exist in current schema
	currentTableInfo, ok := currentSchema[realTableName]
	if !ok {
		// create table with primary key
		query, err := s.addNewTable(logicalDBName, dbType, dbAlias, realTableName, realTableInfo)
		if err != nil {
			return nil, err
		}
		batchedQueries = append(batchedQueries, query)
		currentTableInfo = model.Fields{}
		for realColumnName, realColumnInfo := range realTableInfo {
			temp := model.FieldType{
				FieldName:           realColumnInfo.FieldName,
				IsFieldTypeRequired: realColumnInfo.IsFieldTypeRequired,
				IsList:              realColumnInfo.IsList,
				Kind:                realColumnInfo.Kind,
				IsPrimary:           realColumnInfo.IsPrimary,
				NestedObject:        realColumnInfo.NestedObject,
			}
			currentTableInfo[realColumnName] = &temp
		}
	}

	// Get index maps for each table
	realIndexMap, err := getIndexMap(realTableInfo)
	if err != nil {
		return nil, err
	}
	currentIndexMap, err := getIndexMap(currentTableInfo)
	if err != nil {
		return nil, err
	}

	// Remove the unwanted columns first
	for currentColName, currentColValue := range currentSchema {
		realColValue, ok := realSchema[currentColName]
		// if table doesn't exist handle it grace fully
		if !ok {
			continue
		}
		for currentFieldKey, currentFieldStruct := range currentColValue {
			realField, ok := realColValue[currentFieldKey]
			if !ok || realField.IsLinked {
				// remove field from current table
				c := creationModule{
					dbAlias:           dbAlias,
					logicalDBName:     logicalDBName,
					TableName:         currentColName,
					ColumnName:        currentFieldKey,
					currentColumnInfo: currentFieldStruct,
					currentIndexMap:   currentIndexMap,
					schemaModule:      s,
				}
				if c.currentColumnInfo.IsPrimary {
					return nil, utils.LogError(fmt.Sprintf("Field (%s) with primary key cannot be removed, Delete the table to change primary key", c.ColumnName), "schema", "generateCreationQueries", nil)
				}
				batchedQueries = append(batchedQueries, c.removeColumn(dbType)...)
			}
		}
	}

	for realColumnName, realColumnInfo := range realTableInfo {
		// Ignore the field if its linked. We will be removing the column if it exists later on.
		if realColumnInfo.IsLinked {
			continue
		}
		if err := checkErrors(realColumnInfo); err != nil {
			return nil, err
		}

		// Create the joint table first
		if realColumnInfo.IsForeign {
			if _, p := currentSchema[realColumnInfo.JointTable.Table]; !p {
				if err := s.SchemaCreation(ctx, dbAlias, realColumnInfo.JointTable.Table, logicalDBName, parsedSchema); err != nil {
					return nil, err
				}
			}
		}

		currentColumnInfo, ok := currentTableInfo[realColumnName]
		columnType, err := getSQLType(dbType, realColumnInfo.Kind)
		if err != nil {
			return nil, err
		}
		c := creationModule{
			dbAlias:           dbAlias,
			logicalDBName:     logicalDBName,
			TableName:         realTableName,
			ColumnName:        realColumnName,
			columnType:        columnType,
			currentColumnInfo: currentColumnInfo,
			realColumnInfo:    realColumnInfo,
			currentIndexMap:   currentIndexMap,
			schemaModule:      s,
		}

		if !ok || currentColumnInfo.IsLinked {
			// add field in current table only if its not linked
			if !realColumnInfo.IsLinked {
				queries := c.addColumn(dbType)

				batchedQueries = append(batchedQueries, queries...)
			}

		} else {
			if !realColumnInfo.IsLinked {
				if c.realColumnInfo.Kind != c.currentColumnInfo.Kind {
					// As we are making sure that tables can only be created with primary key, this condition will occur if primary key is removed from a field
					if !c.realColumnInfo.IsPrimary && c.currentColumnInfo.IsPrimary {
						return nil, utils.LogError(fmt.Sprintf(`Cannot change type of field ("%s") primary key exists, Delete the table to change primary key`, c.ColumnName), "schema", "generateCreationQueries", nil)
					}
					// for changing the type of column, drop the column then add new column
					queries := c.modifyColumnType(dbType)
					batchedQueries = append(batchedQueries, queries...)
				}
				if c.currentColumnInfo.IsPrimary && (!c.realColumnInfo.IsPrimary || c.realColumnInfo.IsForeign || !c.realColumnInfo.IsFieldTypeRequired || c.realColumnInfo.IsDefault) {
					return nil, utils.LogError(fmt.Sprintf(`Mutation is not allowed on field ("%s") with primary key, Delete the table to change primary key`, c.ColumnName), "schema", "generateCreationQueries", nil)
				}
				// make changes according to the changes in directives
				queries := c.modifyColumn(dbType)
				batchedQueries = append(batchedQueries, queries...)
			}
		}
	}

	for indexName, currentFields := range currentIndexMap {
		if _, ok := realIndexMap[indexName]; !ok {
			logrus.Println("current index", currentFields, "a", realIndexMap)
			batchedQueries = append(batchedQueries, s.removeIndex(dbType, dbAlias, logicalDBName, tableName, currentFields.IndexName))
		}
	}

	for indexName, fields := range realIndexMap {
		if _, ok := currentIndexMap[indexName]; !ok {
			batchedQueries = append(batchedQueries, s.addIndex(dbType, dbAlias, logicalDBName, tableName, indexName, fields.IsIndexUnique, fields.IndexMap))
			continue
		}
		if !reflect.DeepEqual(fields.IndexMap, cleanIndexMap(currentIndexMap[indexName].IndexMap)) {
			batchedQueries = append(batchedQueries, s.removeIndex(dbType, dbAlias, logicalDBName, tableName, currentIndexMap[indexName].IndexName))
			batchedQueries = append(batchedQueries, s.addIndex(dbType, dbAlias, logicalDBName, tableName, indexName, fields.IsIndexUnique, fields.IndexMap))
		}
	}

	return batchedQueries, nil
}

func cleanIndexMap(v []*model.FieldType) []*model.FieldType {
	for _, fieldType := range v {
		fieldType.IndexInfo.ConstraintName = ""
	}
	return v
}

// SchemaModifyAll modifies all the tables provided
func (s *Schema) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, tables map[string]*config.TableRule) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	crud := config.Crud{}
	crud[dbAlias] = &config.CrudStub{
		Enabled:     true,
		Collections: tables,
	}
	parsedSchema, err := s.Parser(crud)
	if err != nil {
		return err
	}
	for tableName, info := range tables {
		if info.Schema == "" {
			continue
		}
		if err := s.SchemaCreation(ctx, dbAlias, tableName, logicalDBName, parsedSchema); err != nil {
			return err
		}
	}
	return nil
}
