package sql

import (
	"context"

	"github.com/spaceuptech/space-cloud/utils"
)

// DescribeTable return a description of sql table & foreign keys in table
// NOTE: not to be exposed externally
func (s *SQL) DescribeTable(ctx context.Context, project, dbType, col string) ([]utils.FieldType, []utils.ForeignKeysType, error) {
	fields, err := s.getDescribeDetails(ctx, project, dbType, col)
	if err != nil {
		return nil, nil, err
	}
	foreignKeys, err := s.getForeignKeyDetails(ctx, project, dbType, col)
	if err != nil {
		return nil, nil, err
	}
	return fields, foreignKeys, nil
}

func (s *SQL) getDescribeDetails(ctx context.Context, project, dbType, col string) ([]utils.FieldType, error) {
	queryString := ""
	if utils.DBType(dbType) == utils.MySQL {
		queryString = `DESCRIBE ` + project + "." + col
	} else {
		queryString = `SELECT  
		f.attnum AS "Default",  
		f.attnum AS "Extra",
		f.attname AS "Field",  
		pg_catalog.format_type(f.atttypid,f.atttypmod) AS "Type",  
		CASE  
			WHEN p.contype = 'p' THEN 'PRI'  
			WHEN p.contype = 'u' THEN 'UNI'
			ELSE 'f'  
		END AS "Key",
		CASE
			WHEN f.attnotnull = 't' THEN 'YES'
			ELSE 'NO'
		END AS "Null"
	FROM pg_attribute f  
		JOIN pg_class c ON c.oid = f.attrelid  
		JOIN pg_type t ON t.oid = f.atttypid  
		LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = f.attnum  
		LEFT JOIN pg_namespace n ON n.oid = c.relnamespace  
		LEFT JOIN pg_constraint p ON p.conrelid = c.oid AND f.attnum = ANY (p.conkey)  
		LEFT JOIN pg_class AS g ON p.confrelid = g.oid  
	WHERE c.relkind = 'r'::char    
		AND c.relname = '` + col + `'
		AND f.attnum > 0 ORDER BY "Default"`
	}

	rows, err := s.client.Queryx(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.FieldType{}
	for rows.Next() {
		fieldType := new(utils.FieldType)

		if err := rows.StructScan(fieldType); err != nil {
			return nil, err
		}

		result = append(result, *fieldType)
	}
	return result, nil
}

func (s *SQL) getForeignKeyDetails(ctx context.Context, project, dbType, col string) ([]utils.ForeignKeysType, error) {
	queryString := ""
	if utils.DBType(dbType) == utils.MySQL {
		queryString = "select TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = :project and TABLE_NAME = :col"
	} else {
		queryString = `SELECT
		tc.table_name AS "TABLE_NAME", 
		kcu.column_name AS "COLUMN_NAME", 
		tc.constraint_name AS "CONSTRAINT_NAME", 
		ccu.table_name AS "REFERENCED_TABLE_NAME",
		ccu.column_name AS "REFERENCED_COLUMN_NAME"
	FROM 
		information_schema.table_constraints AS tc 
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		  AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
		  ON ccu.constraint_name = tc.constraint_name
		  AND ccu.table_schema = tc.table_schema
	WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name= :col
	`
	}
	rows, err := s.client.NamedQuery(queryString, map[string]interface{}{"col": col, "project": project})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.ForeignKeysType{}
	for rows.Next() {
		foreignKey := new(utils.ForeignKeysType)

		if err := rows.StructScan(foreignKey); err != nil {
			return nil, err
		}

		result = append(result, *foreignKey)
	}
	return result, nil
}
