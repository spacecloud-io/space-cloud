package schema

import (
	"context"
	"errors"
	"reflect"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type creationModule struct {
	dbAlias, project, TableName, ColumnName, columnType string
	currentColumnInfo, realColumnInfo                   *SchemaFieldType
	schemaModule                                        *Schema
	removeProjectScope                                  bool
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbAlias, tableName, project string, parsedSchema schemaType) error {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return err
	}

	// Return gracefully if db type is mongo
	if dbType == string(utils.Mongo) || dbType == string(utils.BoltDB) {
		return nil
	}

	if err := s.crud.CreateProjectIfNotExists(ctx, project, dbAlias); err != nil {
		return err
	}

	currentSchema, _ := s.Inspector(ctx, dbAlias, project, tableName)

	queries, err := s.generateCreationQueries(ctx, dbAlias, tableName, project, parsedSchema, currentSchema)
	if err != nil {
		return err
	}
	return s.crud.RawBatch(ctx, dbAlias, queries)
}

func (s *Schema) generateCreationQueries(ctx context.Context, dbAlias, tableName, project string, parsedSchema schemaType, currentSchema schemaCollection) ([]string, error) {
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
		query, err := addNewTable(project, dbAlias, realTableName, realTableInfo, s.removeProjectScope)
		if err != nil {
			return nil, err
		}
		batchedQueries = append(batchedQueries, query)
		currentTableInfo = SchemaFields{}
		for realColumnName, realColumnInfo := range realTableInfo {
			temp := SchemaFieldType{
				FieldName:           realColumnInfo.FieldName,
				IsFieldTypeRequired: realColumnInfo.IsFieldTypeRequired,
				IsList:              realColumnInfo.IsList,
				Kind:                realColumnInfo.Kind,
				IsPrimary:           realColumnInfo.IsPrimary,
				nestedObject:        realColumnInfo.nestedObject,
			}
			currentTableInfo[realColumnName] = &temp
		}
	}

	for realColumnName, realColumnInfo := range realTableInfo {
		// Ignore the field if its linked
		if realColumnInfo.IsLinked {
			continue
		}
		if err := checkErrors(realColumnInfo); err != nil {
			return nil, err
		}

		// Create the joint table first
		if realColumnInfo.IsForeign {
			if _, p := currentSchema[realColumnInfo.JointTable.Table]; !p {
				if err := s.SchemaCreation(ctx, dbAlias, realColumnInfo.JointTable.Table, project, parsedSchema); err != nil {
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
			dbAlias:            dbAlias,
			project:            project,
			TableName:          realTableName,
			ColumnName:         realColumnName,
			columnType:         columnType,
			currentColumnInfo:  currentColumnInfo,
			realColumnInfo:     realColumnInfo,
			schemaModule:       s,
			removeProjectScope: s.removeProjectScope,
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
					// for changing the type of column, drop the column then add new column
					queries := c.modifyColumnType(dbType)

					batchedQueries = append(batchedQueries, queries...)
				} else {
					// make changes according to the changes in directives
					queries := c.modifyColumn()

					batchedQueries = append(batchedQueries, queries...)
				}
			}
		}
	}

	for currentColName, currentColValue := range currentSchema {
		realColValue, ok := realSchema[currentColName]
		// if table doesn't exist handle it grace fully
		if !ok {
			continue
		}
		for currentFieldKey, currentFieldStruct := range currentColValue {
			realField, ok := realColValue[currentFieldKey]
			if !ok || realField.IsLinked {
				// remove field from current tabel
				c := creationModule{
					dbAlias:            dbAlias,
					project:            project,
					TableName:          currentColName,
					ColumnName:         currentFieldKey,
					currentColumnInfo:  currentFieldStruct,
					removeProjectScope: s.removeProjectScope,
				}
				if c.currentColumnInfo.IsForeign {
					batchedQueries = append(batchedQueries, c.removeForeignKey()...)
				}
				batchedQueries = append(batchedQueries, c.removeColumn())
			}
		}
	}

	realIndexMap, err := getRealIndexMap(realTableInfo)
	if err != nil {
		return nil, err
	}
	currentIndexMap, err := getCurrentIndexMap(currentTableInfo)
	if err != nil {
		return nil, err
	}
	for indexName, fields := range realIndexMap {
		if _, ok := currentIndexMap[indexName]; !ok {
			batchedQueries = append(batchedQueries, addIndex(dbType, project, tableName, indexName, fields.IsIndexUnique, s.removeProjectScope, fields.IndexMap))
			continue
		}
		if !reflect.DeepEqual(fields.IndexMap, currentIndexMap[indexName].IndexMap) {
			batchedQueries = append(batchedQueries, removeIndex(dbType, project, tableName, indexName, s.removeProjectScope))
			batchedQueries = append(batchedQueries, addIndex(dbType, project, tableName, indexName, fields.IsIndexUnique, s.removeProjectScope, fields.IndexMap))
		}
	}
	for indexName, _ := range currentIndexMap {
		if _, ok := realIndexMap[indexName]; !ok {
			batchedQueries = append(batchedQueries, removeIndex(dbType, project, tableName, indexName, s.removeProjectScope))
		}
	}

	return batchedQueries, nil
}

// SchemaModifyAll modifies all the tables provided
func (s *Schema) SchemaModifyAll(ctx context.Context, dbAlias, project string, tables map[string]*config.TableRule) error {
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
		if err := s.SchemaCreation(ctx, dbAlias, tableName, project, parsedSchema); err != nil {
			return err
		}
	}
	return nil
}
