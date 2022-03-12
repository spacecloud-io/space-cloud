package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-test/deep"
	"github.com/spaceuptech/helpers"

	"github.com/spacecloud-io/space-cloud/model"
)

// PrepareCreationQueries generates a set of SQL queries that need to be run in order to apply the changes made in the schema
// to the database.
func PrepareCreationQueries(ctx context.Context, dbType, tableName, dbName string, newSchema, currentSchema model.CollectionSchemas, fn createSchemaFunc) ([]string, error) {
	// Return nil, if no tables are present in schema
	if len(newSchema) == 0 {
		return nil, nil
	}

	batchedQueries := []string{}

	realTableName := tableName
	realTableInfo, p1 := newSchema[realTableName]
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
		query, err := addNewTable(ctx, dbName, dbType, realTableName, realTableInfo)
		if err != nil {
			return nil, err
		}
		batchedQueries = append(batchedQueries, query)
		currentTableInfo = model.FieldSchemas{}
		for realColumnName, realColumnInfo := range realTableInfo {
			temp := model.FieldType{
				TypeIDSize:          realColumnInfo.TypeIDSize,
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
	realIndexMap, err := getIndexMap(ctx, realTableInfo)
	if err != nil {
		return nil, err
	}
	currentIndexMap, err := getIndexMap(ctx, currentTableInfo)
	if err != nil {
		return nil, err
	}

	// Remove the unwanted columns first
	for currentColName, currentColValue := range currentSchema {
		realColValue, ok := newSchema[currentColName]
		// if table doesn't exist handle it grace fully
		if !ok {
			continue
		}

		for currentFieldKey, currentFieldStruct := range currentColValue {
			if currentFieldStruct.IsLinked {
				continue
			}
			realField, ok := realColValue[currentFieldKey]
			if !ok || realField.IsLinked {
				// remove field from current table
				c := creationModule{
					dbType:            dbType,
					dbName:            dbName,
					TableName:         currentColName,
					ColumnName:        currentFieldKey,
					currentColumnInfo: currentFieldStruct,
					currentIndexMap:   currentIndexMap,
				}
				if c.currentColumnInfo.IsPrimary {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Field (%s) with primary key cannot be removed, Delete the table to change primary key", c.ColumnName), nil, nil)
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
		if err := checkErrors(ctx, realColumnInfo); err != nil {
			return nil, err
		}

		// Create the joint table first
		if realColumnInfo.IsForeign {
			if _, p := currentSchema[realColumnInfo.JointTable.Table]; !p {
				if err := fn(ctx, realColumnInfo.JointTable.Table, newSchema); err != nil {
					return nil, err
				}
			}
		}

		currentColumnInfo, ok := currentTableInfo[realColumnName]
		columnType, err := getSQLType(ctx, dbType, realColumnInfo)
		if err != nil {
			return nil, err
		}
		c := creationModule{
			dbType:            dbType,
			dbName:            dbName,
			TableName:         realTableName,
			ColumnName:        realColumnName,
			columnType:        columnType,
			currentColumnInfo: currentColumnInfo,
			realColumnInfo:    realColumnInfo,
			currentIndexMap:   currentIndexMap,
		}

		if !ok || currentColumnInfo.IsLinked {
			// add field in current table only if its not linked
			if !realColumnInfo.IsLinked {
				queries := c.addColumn(dbType)

				batchedQueries = append(batchedQueries, queries...)
			}

		} else {
			if !realColumnInfo.IsLinked {
				if arr := deep.Equal(c.realColumnInfo.Args, c.currentColumnInfo.Args); c.realColumnInfo.Kind != c.currentColumnInfo.Kind || (c.realColumnInfo.TypeIDSize != c.currentColumnInfo.TypeIDSize) || (c.currentColumnInfo.Args != nil && len(arr) > 0) {
					// As we are making sure that tables can only be created with primary key, this condition will occur if primary key is removed from a field
					if c.realColumnInfo.IsPrimary {
						return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf(`Cannot change type of field ("%s") primary key exists, Delete the table to change primary key`, c.ColumnName), nil, nil)
					}
					if !c.realColumnInfo.IsPrimary && c.currentColumnInfo.IsPrimary {
						return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf(`Cannot remove primary constraint on field ("%s") primary key exists, Delete the table to change primary key`, c.ColumnName), nil, nil)
					}
					// for changing the type of column, drop the column then add new column
					queries := c.modifyColumnType(dbType)
					batchedQueries = append(batchedQueries, queries...)
				}
				if c.currentColumnInfo.IsPrimary && (!c.realColumnInfo.IsPrimary || c.realColumnInfo.IsForeign || !c.realColumnInfo.IsFieldTypeRequired || c.realColumnInfo.IsDefault) {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf(`Mutation is not allowed on field ("%s") with primary key, Delete the table to change primary key`, c.ColumnName), nil, nil)
				}
				// make changes according to the changes in directives
				queries := c.modifyColumn(dbType)
				batchedQueries = append(batchedQueries, queries...)
			}
		}
	}

	for indexName, currentFields := range currentIndexMap {
		if _, ok := realIndexMap[indexName]; !ok {
			batchedQueries = append(batchedQueries, removeIndex(dbType, dbName, tableName, currentFields.IndexName))
		}
	}

	for indexName, fields := range realIndexMap {
		if _, ok := currentIndexMap[indexName]; !ok {
			batchedQueries = append(batchedQueries, addIndex(dbType, dbName, tableName, indexName, fields.IsIndexUnique, fields.IndexTableProperties))
			continue
		}
		if arr := deep.Equal(fields.IndexTableProperties, cleanIndexMap(currentIndexMap[indexName].IndexTableProperties)); len(arr) > 0 {
			batchedQueries = append(batchedQueries, removeIndex(dbType, dbName, tableName, currentIndexMap[indexName].IndexName))
			batchedQueries = append(batchedQueries, addIndex(dbType, dbName, tableName, indexName, fields.IsIndexUnique, fields.IndexTableProperties))
		}
	}

	return batchedQueries, nil
}
