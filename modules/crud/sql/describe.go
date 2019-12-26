package sql

import (
	"context"
	"errors"

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
	args := []interface{}{}
	switch utils.DBType(dbType) {
	case utils.MySQL:
		queryString = `DESCRIBE ` + project + "." + col
	case utils.Postgres:
		queryString = `SELECT  
		adsrc AS "Default",  
		f.attnum AS "Extra",
		f.attname AS "Field",  
		pg_catalog.format_type(f.atttypid,f.atttypmod) AS "Type",  
		CASE  
			WHEN p.contype = 'p' THEN 'PRI'  
			WHEN p.contype = 'u' THEN 'UNI'
			ELSE 'f'  
		END AS "Key",
		CASE
			WHEN f.attnotnull = 't' THEN 'NO'
			ELSE 'YES'
		END AS "Null"
FROM pg_class pclass,pg_attribute f
		JOIN pg_class c ON c.oid = f.attrelid  
		JOIN pg_type t ON t.oid = f.atttypid  
		LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = f.attnum  
		LEFT JOIN pg_namespace n ON n.oid = c.relnamespace  
		LEFT JOIN pg_constraint p ON p.conrelid = c.oid AND f.attnum = ANY (p.conkey)  
		LEFT JOIN pg_class AS g ON p.confrelid = g.oid  
	WHERE c.relkind = 'r'::char    
	  AND pclass.relname = $1
		AND c.relname = $1
		AND n.nspname = $2
		AND f.attnum > 0 ORDER BY "Default"`

		args = append(args, col, project)
	case utils.SqlServer:

		queryString = `SELECT DISTINCT C.COLUMN_NAME as 'Field', C.IS_NULLABLE as 'Null' , C.DATA_TYPE as 'Type',C.COLUMN_DEFAULT as 'Default',C.DATA_TYPE as 'Extra',
       CASE
           WHEN TC.CONSTRAINT_TYPE = 'PRIMARY KEY' THEN 'PRI'
           WHEN TC.CONSTRAINT_TYPE = 'UNIQUE' THEN 'UNI'
           ELSE isnull(TC.CONSTRAINT_TYPE,'NULL')
           END AS 'Key'
FROM INFORMATION_SCHEMA.COLUMNS AS C
         FULL JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE AS CC
                   ON C.COLUMN_NAME = CC.COLUMN_NAME
         FULL JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
                   ON CC.CONSTRAINT_NAME = TC.CONSTRAINT_NAME
WHERE C.TABLE_SCHEMA=@p2 AND C.table_name = @p1`

		args = append(args, col, project)
	}
	rows, err := s.client.QueryxContext(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []utils.FieldType{}
	count := 0
	for rows.Next() {
		count++
		fieldType := new(utils.FieldType)

		if err := rows.StructScan(fieldType); err != nil {
			return nil, err
		}

		result = append(result, *fieldType)
	}
	if count == 0 {
		return result, errors.New(dbType + ":" + col + " not found during inspection")
	}
	return result, nil
}

func (s *SQL) getForeignKeyDetails(ctx context.Context, project, dbType, col string) ([]utils.ForeignKeysType, error) {
	queryString := ""
	switch utils.DBType(dbType) {

	case utils.MySQL:
		queryString = "select TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = ? and TABLE_NAME = ?"
	case utils.Postgres:
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
	WHERE tc.constraint_type = 'FOREIGN KEY'  AND tc.table_schema = $1  AND tc.table_name= $2
	`
	case utils.SqlServer:
		queryString = `SELECT 
		CCU.TABLE_NAME, CCU.COLUMN_NAME, CCU.CONSTRAINT_NAME,
		isnull(KCU.TABLE_NAME,'') AS 'REFERENCED_TABLE_NAME', isnull(KCU.COLUMN_NAME,'') AS 'REFERENCED_COLUMN_NAME'
	FROM INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE CCU
		FULL JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
			ON CCU.CONSTRAINT_NAME = RC.CONSTRAINT_NAME 
		FULL JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE KCU 
			ON KCU.CONSTRAINT_NAME = RC.UNIQUE_CONSTRAINT_NAME  
	WHERE CCU.TABLE_SCHEMA = @p1 AND CCU.TABLE_NAME= @p2`
	}
	rows, err := s.client.QueryxContext(ctx, queryString, []interface{}{project, col}...)
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
