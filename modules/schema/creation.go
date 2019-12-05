package schema

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

type creationModule struct {
	dbType, project, TableName, ColumnName, columnType string
	currentColumnInfo, realColumnInfo                  *SchemaFieldType
	schemaModule                                       *Schema
	removeProjectScope                                 bool
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbType, tableName, project string, parsedSchema schemaType) error {
	// Return gracefully if db type is mongo
	if dbType == string(utils.Mongo) {
		return nil
	}

	// Return if no tables are present in schema
	if len(parsedSchema[dbType]) == 0 {
		return nil
	}

	if err := s.crud.CreateProjectIfNotExists(ctx, project, dbType); err != nil {
		return err
	}

	currentSchema, _ := s.Inspector(ctx, dbType, project, tableName)
	realSchema := parsedSchema[dbType]
	batchedQueries := []string{}

	realTableName := tableName
	realTableInfo, p1 := realSchema[realTableName]
	if !p1 {
		if _, p2 := currentSchema[realTableName]; p2 {
			return nil
		}

		return errors.New("Schema not provided for table: " + tableName)
	}

	// check if table exist in current schema

	currentTableInfo, ok := currentSchema[realTableName]
	if !ok {
		// create table with primary key
		query, err := addNewTable(project, dbType, realTableName, realTableInfo, s.removeProjectScope)
		if err != nil {
			return err
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
			return err
		}

		// Create the joint table first
		if realColumnInfo.IsForeign {
			if err := s.SchemaCreation(ctx, dbType, realColumnInfo.JointTable.Table, project, parsedSchema); err != nil {
				return err
			}
		}

		currentFieldStruct, ok := currentTableInfo[realColumnName]
		columnType, err := getSQLType(dbType, realColumnInfo.Kind)
		if err != nil {
			return err
		}
		c := creationModule{
			dbType:             dbType,
			project:            project,
			TableName:          realTableName,
			ColumnName:         realColumnName,
			columnType:         columnType,
			currentColumnInfo:  currentFieldStruct,
			realColumnInfo:     realColumnInfo,
			schemaModule:       s,
			removeProjectScope: s.removeProjectScope,
		}

		if !ok {
			// add field in current table only if its not linked
			if !realColumnInfo.IsLinked {
				queries, err := c.addField(ctx)
				if err != nil {
					return err
				}

				batchedQueries = append(batchedQueries, queries...)
			}

		} else {
			// modify removing then adding
			if !realColumnInfo.IsLinked {
				queries, err := c.modifyField(ctx)
				batchedQueries = append(batchedQueries, queries...)
				if err != nil {
					return err
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
					dbType:             dbType,
					project:            project,
					TableName:          currentColName,
					ColumnName:         currentFieldKey,
					currentColumnInfo:  currentFieldStruct,
					removeProjectScope: s.removeProjectScope,
				}

				if c.currentColumnInfo.IsForeign {
					batchedQueries = append(batchedQueries, c.removeForeignKey()...)
				}

				batchedQueries = append(batchedQueries, c.removeField())
			}
		}
	}
	return s.crud.RawBatch(ctx, dbType, batchedQueries)
}

// SchemaModifyAll modifies all the tables provided
func (s *Schema) SchemaModifyAll(ctx context.Context, dbType, project string, tables map[string]*config.TableRule) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	crud := config.Crud{}
	crud[dbType] = &config.CrudStub{
		Enabled:     true,
		Collections: tables,
	}

	parsedSchema, err := s.parser(crud)
	if err != nil {
		return err
	}

	for tableName, info := range tables {
		if info.Schema == "" {
			continue
		}

		if err := s.SchemaCreation(ctx, dbType, tableName, project, parsedSchema); err != nil {
			return err
		}
	}
	return nil
}
