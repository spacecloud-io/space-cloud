package schema

import (
	"context"
	"encoding/json"
	"fmt"

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

	parsedSchema, err := s.parser(crud)
	if err != nil {
		return err
	}

	currentSchema, err := s.Inspector(ctx, dbType, project, col)
	if err != nil {
		return err
	}

	realSchema := parsedSchema[dbType]

	batchedQueries := []string{}

	for realColName, realColValue := range realSchema {
		// check if table exist in current schema
		currentColValue := currentSchema[realColName]

		for realFieldKey, realFieldStruct := range realColValue {
			if err := checkErrors(realFieldStruct); err != nil {
				return err
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
		realColValue := realSchema[currentColName]

		for currentFieldKey := range currentColValue {
			c := creationModule{
				project:  project,
				ColName:  currentColName,
				FieldKey: currentFieldKey,
			}
			_, ok := realColValue[currentFieldKey]
			if !ok {
				// remove field from current tabel
				batchedQueries = append(batchedQueries, c.removeField())
			}
		}
	}
	b, err := json.MarshalIndent(batchedQueries, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
	return s.crud.RawBatch(ctx, dbType, batchedQueries)
}

func (c *creationModule) addField(ctx context.Context) ([]string, error) {
	var queryString string
	var queries []string

	if c.columnType != "" {
		// add a new column with data type as columntype
		queryString = "ALTER TABLE " + c.project + "." + c.ColName + " ADD " + c.FieldKey + " " + c.columnType
		queries = append(queries, queryString)
	}

	if c.realFieldStruct.IsFieldTypeRequired {
		// make the new column not null
		queryString = "ALTER TABLE " + c.project + "." + c.ColName + " MODIFY " + c.FieldKey + " " + c.columnType + " NOT NULL"
		queries = append(queries, queryString)
	}
	queries = append(queries, c.addDirective(ctx)...)
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
		queries = append(queries, c.modifyNotNull())
	}

	if c.realFieldStruct.Directive != c.currentFieldStruct.Directive {
		queries = append(queries, c.removeDirective()...)
		queries = append(queries, c.addDirective(ctx)...)
	}
	return queries, nil
}

func (c *creationModule) addDirective(ctx context.Context) []string {
	queries := []string{}
	switch c.realFieldStruct.Directive {
	case directiveId:
		return append(queries, c.addPrimaryKey())
	case directiveUnique:
		return append(queries, c.addUniqueKey())
	case directiveRelation:
		nestedSchema, err := c.schemaModule.Inspector(ctx, c.dbType, c.project, c.realFieldStruct.JointTable.TableName)
		if err != nil {
			return nil
		}
		value, ok := nestedSchema[c.realFieldStruct.JointTable.TableName]
		if !ok {
			// error handling
		}
		val, ok := value[c.realFieldStruct.JointTable.TableField]
		if !ok {
			// error handling
		}
		colType, err := getSQLType(val.Kind)
		if err != nil {
			// error handlin
		}
		if colType == typeObject || colType == typeJoin || val.IsList {
			// err
		}
		temp := creationModule{
			dbType:     c.dbType,
			project:    c.project,
			ColName:    c.ColName,
			FieldKey:   c.FieldKey,
			columnType: colType,
		}
		return append(queries, temp.addNewColumn(), c.addForeignKey())
	}
	return queries
}

func (c *creationModule) removeDirective() []string {
	queries := []string{}
	switch c.currentFieldStruct.Directive {
	case directiveId:
		return append(queries, c.removePrimaryKey())
	case directiveUnique:
		return append(queries, c.removeUniqueKey())
	case directiveRelation:
		return append(queries, c.removeForeignKey())
	}
	return queries
}
