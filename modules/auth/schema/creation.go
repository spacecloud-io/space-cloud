package schema

import (
	"context"

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
			if err := s.SchemaCreation(ctx, dbName, colName, s.project, tableRule.Schema); err != nil {
				return err
			}
		}
	}

	return nil
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbType, col, project, schema string) error {
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

	for realColName, realColValue := range realSchema {
		// check if table exist in current schema
		currentColValue, ok := currentSchema[realColName]
		if !ok {
			// create table with primary key
			query, err := addNewTable(project, dbType, realColName, realColValue)
			if err != nil {
				return err
			}

			batchedQueries = append(batchedQueries, query)
			if err := s.crud.RawBatch(ctx, dbType, batchedQueries); err != nil {
				return err
			}
			return s.SchemaCreation(ctx, dbType, col, project, schema)
		}
		for realFieldKey, realFieldStruct := range realColValue {
			if err := checkErrors(realFieldStruct); err != nil {
				return err
			}
			if realFieldStruct.IsList && (realFieldStruct.Directive == directiveRelation) { // as directive is relation for array type don't generate queries
				continue
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
			}

			if !ok {
				// add field in current table

				queries, err := c.addField(ctx)
				batchedQueries = append(batchedQueries, queries...)
				if err != nil {
					return err
				}

			} else {
				// modify removing then adding
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

	return s.crud.RawBatch(ctx, dbType, batchedQueries)
}

func (s *Schema) SchemaJoin(ctx context.Context, parsedColValue schemaCollection, dbType, colName, project string, v config.CrudStub) error {
	temp := &schemaFieldType{}
	for _, fieldValue := range parsedColValue[colName] {
		temp = fieldValue
		if fieldValue.Kind == typeJoin && fieldValue.Directive == directiveRelation && !fieldValue.IsList {
			if err := s.SchemaJoin(ctx, parsedColValue, dbType, fieldValue.JointTable.TableName, project, v); err != nil {
				return err
			}
		}
	}
	return s.SchemaCreation(ctx, dbType, temp.JointTable.TableName, project, v.Collections[temp.JointTable.TableName].Schema)
}
