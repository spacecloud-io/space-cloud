package schema

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

type creationModule struct {
	dbType, project, ColName, FieldKey, columnType string
	currentFieldStruct, realFieldStruct            *SchemaFieldType
	schemaModule                                   *Schema
	removeProjectScope                             bool
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbType, col, project string, parsedSchema schemaType) error {

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

	currentSchema, _ := s.Inspector(ctx, dbType, project, col)

	realSchema := parsedSchema[dbType]
	batchedQueries := []string{}

	realColName := col
	realColValue, p1 := realSchema[realColName]
	if !p1 {
		if _, p2 := currentSchema[realColName]; p2 {
			return nil
		}

		return errors.New("Schema not provided for table: " + col)
	}

	// check if table exist in current schema
	currentColValue, ok := currentSchema[realColName]
	if !ok {
		// create table with primary key
		query, err := addNewTable(project, dbType, realColName, realColValue, s.removeProjectScope)
		if err != nil {
			return err
		}

		batchedQueries = append(batchedQueries, query)

		currentColValue = SchemaFields{}
		for realFieldKey, realFieldStruct := range realColValue {
			temp := SchemaFieldType{
				FieldName:           realFieldStruct.FieldName,
				IsFieldTypeRequired: realFieldStruct.IsFieldTypeRequired,
				IsList:              realFieldStruct.IsList,
				Kind:                realFieldStruct.Kind,
				IsPrimary:           realFieldStruct.IsPrimary,
				nestedObject:        realFieldStruct.nestedObject,
			}

			currentColValue[realFieldKey] = &temp
		}
	}
	for realFieldKey, realFieldStruct := range realColValue {
		// Ignore the field if its linked
		if realFieldStruct.IsLinked {
			continue
		}
		if err := checkErrors(realFieldStruct); err != nil {
			return err
		}

		// Create the joint table first
		if realFieldStruct.IsForeign {
			if err := s.SchemaCreation(ctx, dbType, realFieldStruct.JointTable.Table, project, parsedSchema); err != nil {
				return err
			}
		}

		currentFieldStruct, ok := currentColValue[realFieldKey]
		columnType, err := getSQLType(dbType, realFieldStruct.Kind)
		if err != nil {
			return err
		}
		c := creationModule{
			dbType:             dbType,
			project:            project,
			ColName:            realColName,
			FieldKey:           realFieldKey,
			columnType:         columnType,
			currentFieldStruct: currentFieldStruct,
			realFieldStruct:    realFieldStruct,
			schemaModule:       s,
			removeProjectScope: s.removeProjectScope,
		}

		if !ok {
			// add field in current table only if its not linked
			if !realFieldStruct.IsLinked {
				queries, err := c.addField(ctx)
				if err != nil {
					return err
				}

				batchedQueries = append(batchedQueries, queries...)
			}

		} else {
			// modify removing then adding
			if !realFieldStruct.IsLinked {
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
					ColName:            currentColName,
					FieldKey:           currentFieldKey,
					currentFieldStruct: currentFieldStruct,
					removeProjectScope: s.removeProjectScope,
				}

				if c.currentFieldStruct.IsForeign {
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
