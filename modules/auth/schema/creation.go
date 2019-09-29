package schema

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/config"
)

type creationModule struct {
	dbType, project, ColName, FieldKey, columnType string
	currentFieldStruct, realFieldStruct            *schemaFieldType
	schemaModule                                   *Schema
}

// SchemaCreation creates or alters tables of sql
func (s *Schema) SchemaCreation(ctx context.Context, dbType, col, project, schema string) error {
	crudCol := map[string]*config.TableRule{}
	crudCol[col] = &config.TableRule{
		Schema: schema,
	}

	crud := config.Crud{}
	crud[dbType] = &config.CrudStub{
		Enabled:     true,
		Collections: crudCol,
	}

	if err := s.crud.CreateProjectIfNotExists(ctx, project, dbType); err != nil {
		return err
	}

	parsedSchema, err := s.parser(crud)
	if err != nil {
		return nil
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

		for currentFieldKey := range currentColValue {
			_, ok := realColValue[currentFieldKey]
			if !ok {
				// remove field from current tabel
				c := creationModule{
					project:  project,
					ColName:  currentColName,
					FieldKey: currentFieldKey,
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
		} else {
			queries = append(queries, c.removeDirective()...)
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
		// get schema of referenced table & check if referenced column exist
		nestedSchema, err := c.schemaModule.Inspector(ctx, c.dbType, c.project, c.realFieldStruct.JointTable.TableName)
		if err != nil {
			return nil, err
		}
		value, ok := nestedSchema[c.realFieldStruct.JointTable.TableName]
		if !ok {
			return nil, errors.New("schema creation: foreign key referenced table not found")
		}
		val, ok := value[c.realFieldStruct.JointTable.TableField]
		if !ok {
			return nil, errors.New("schema creation: field name not found in referenced table for foreign keys")
		}
		colType, err := getSQLType(val.Kind)
		if err != nil {
			return nil, err
		}
		if colType == typeObject || colType == typeJoin || val.IsList {
			return nil, errors.New("schema creation: found incorrect type or array in foreign key")
		}
		temp := creationModule{
			dbType:     c.dbType,
			project:    c.project,
			ColName:    c.ColName,
			FieldKey:   c.FieldKey,
			columnType: colType,
		}
		queries = append(queries, temp.modifyColumnType(), c.addForeignKey())
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
