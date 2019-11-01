package schema

import (
	"context"
	"errors"
	"log"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

type creationModule struct {
	dbType, project, ColName, FieldKey, columnType string
	currentFieldStruct, realFieldStruct            *schemaFieldType
	schemaModule                                   *Schema
}

func (s *Schema) ModifyAllCollections(ctx context.Context, conf config.Crud) error {
	for dbName, crudStubValue := range conf {
		if utils.DBType(dbName) == utils.Mongo {
			continue
		}
		if !crudStubValue.Enabled {
			continue
		}

		for colName, tableRule := range crudStubValue.Collections {
			if _, err := s.SchemaCreation(ctx, dbName, colName, s.project, tableRule.Schema, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbType, col, project, schema, skipTable string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	crudCol := map[string]*config.TableRule{}
	crudCol[col] = &config.TableRule{
		Schema: schema,
	}

	crud := config.Crud{}
	crud[dbType] = &config.CrudStub{
		Enabled:     true,
		Collections: crudCol,
	}

	parsedSchema, err := s.parser(crud)
	if err != nil {
		return "", nil
	}

	// Return if no tables are present in schema
	if len(parsedSchema[dbType]) == 0 {
		return "", nil

	}

	if err := s.crud.CreateProjectIfNotExists(ctx, project, dbType); err != nil {
		return "", nil
	}

	currentSchema, _ := s.Inspector(ctx, dbType, project, col)

	realSchema := parsedSchema[dbType]
	batchedQueries := []string{}

	for realColName, realColValue := range realSchema {
		// check if table exist in current schema
		currentColValue, ok := currentSchema[realColName]
		if !ok {
			// create table with primary key
			query, err := addNewTable(project, dbType, realColName, realColValue)
			if err != nil {
				return "", nil
			}

			batchedQueries = append(batchedQueries, query)
			if err := s.crud.RawBatch(ctx, dbType, batchedQueries); err != nil {
				return "", nil
			}
			return s.SchemaCreation(ctx, dbType, col, project, schema, "")
		}
		for realFieldKey, realFieldStruct := range realColValue {
			if err := checkErrors(realFieldStruct); err != nil {
				return "", nil
			}
			if realFieldStruct.IsList && (realFieldStruct.Directive == directiveRelation) { // as directive is relation for array type don't generate queries
				continue
			}
			if !realFieldStruct.IsList && realFieldStruct.Directive == directiveRelation && realFieldStruct.Kind == typeJoin && realFieldStruct.JointTable.TableName != skipTable {
				return realFieldStruct.JointTable.TableName, nil
			}
			currentFieldStruct, ok := currentColValue[realFieldKey]
			columnType, err := getSQLType(dbType, realFieldStruct.Kind)
			if err != nil {
				return "", nil
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
			}

			if !ok {
				// add field in current table

				queries, err := c.addField(ctx)
				batchedQueries = append(batchedQueries, queries...)
				if err != nil {
					return "", nil
				}

			} else {
				// modify removing then adding
				queries, err := c.modifyField(ctx)
				batchedQueries = append(batchedQueries, queries...)
				if err != nil {
					return "", nil
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
			_, ok := realColValue[currentFieldKey]
			if !ok {
				// remove field from current tabel
				c := creationModule{
					dbType:             dbType,
					project:            project,
					ColName:            currentColName,
					FieldKey:           currentFieldKey,
					currentFieldStruct: currentFieldStruct,
				}

				if c.currentFieldStruct.Directive == directiveRelation {
					batchedQueries = append(batchedQueries, c.removeForeignKey()...)
				}

				batchedQueries = append(batchedQueries, c.removeField())
			}
		}
	}

	return "", s.crud.RawBatch(ctx, dbType, batchedQueries)
}

func (s *Schema) SchemaCreationWithObject(ctx context.Context, dbType, project string, tables map[string]*config.TableRule) error {
	tablesDone := map[string]bool{}
	var err error
	for tableName, tableInfo := range tables {
		log.Println("called creation object", tablesDone)
		_, ok := tablesDone[tableName]
		if !ok {
			if err = s.recursiveSchemaCreation(ctx, dbType, project, tableName, tableInfo.Schema, tables, tablesDone); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Schema) recursiveSchemaCreation(ctx context.Context, dbType, project, colName, schema string, collectionMap map[string]*config.TableRule, tablesDone map[string]bool) error {
	skipTable := ""
	for {
		tableName, err := s.SchemaCreation(ctx, dbType, colName, project, schema, skipTable)
		if err != nil {
			return err
		}
		_, ok := tablesDone[tableName]
		if ok {
			skipTable = tableName
			continue
		}
		if tableName != "" && !ok {
			collection, ok := collectionMap[tableName]
			if !ok {
				return errors.New("Schema not provided for " + tableName)
			}
			if err = s.recursiveSchemaCreation(ctx, dbType, project, tableName, collection.Schema, collectionMap, tablesDone); err != nil {
				return err
			}
			skipTable = tableName
		} else {
			break
		}
	}
	tablesDone[colName] = true
	return nil
}
