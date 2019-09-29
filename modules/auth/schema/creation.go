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
			query, err := addNewTable(project, realColName, realColValue)
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
			columnType, err := getSQLType(realFieldStruct.Kind)
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

// func (s *Schema) createProject(dbType string) {
// 	s.crud.
// }
