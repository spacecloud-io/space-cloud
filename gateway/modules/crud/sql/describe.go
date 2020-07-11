package sql

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// DescribeTable return a description of sql table & foreign keys in table
// NOTE: not to be exposed externally
func (s *SQL) DescribeTable(ctx context.Context, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error) {
	fields, err := s.getDescribeDetails(ctx, s.name, col)
	if err != nil {
		return nil, nil, nil, err
	}
	foreignKeys, err := s.getForeignKeyDetails(ctx, s.name, col)
	if err != nil {
		return nil, nil, nil, err
	}
	index, err := s.getIndexDetails(ctx, s.name, col)
	if err != nil {
		return nil, nil, nil, err
	}

	return fields, foreignKeys, index, nil
}

func (s *SQL) getDescribeDetails(ctx context.Context, project, col string) ([]utils.FieldType, error) {
	queryString := ""
	args := []interface{}{}
	switch utils.DBType(s.dbType) {
	case utils.MySQL:
		queryString = `select column_name as 'Field',is_nullable as 'Null',column_key as 'Key',
case when data_type = 'varchar' then concat(DATA_TYPE,'(',CHARACTER_MAXIMUM_LENGTH,')') else DATA_TYPE end as 'Type',
CASE 
	WHEN column_default = '1' THEN 'true'
	WHEN column_default = '0' THEN 'false'
	ELSE coalesce(column_default,'')
END AS 'Default'
from information_schema.columns
where (table_name,table_schema) = (?,?);`
		args = append(args, col, project)

	case utils.Postgres:
		queryString = `SELECT isc.column_name AS "Field", SPLIT_PART(REPLACE(coalesce(column_default,''),'''',''), '::', 1) AS "Default" ,isc.data_type AS "Type",isc.is_nullable AS "Null",
CASE
    WHEN t.constraint_type = 'PRIMARY KEY' THEN 'PRI'
    WHEN t.constraint_type = 'UNIQUE' THEN 'UNI'
    ELSE ''
END AS "Key"
FROM information_schema.columns isc
    left join (select cu.table_schema, cu.table_name, cu.column_name, istc.constraint_type 
    	from information_schema.constraint_column_usage cu 
    	left join information_schema.table_constraints istc on (istc.table_schema,istc.table_name, istc.constraint_name) = (cu.table_schema,cu.table_name, cu.constraint_name) 
    	where istc.constraint_type != 'CHECK') t
    on (t.table_schema, t.table_name, t.column_name) = (isc.table_schema, isc.table_name, isc.column_name)
WHERE (isc.table_schema, isc.table_name) = ($2, $1)
ORDER BY isc.ordinal_position;`

		args = append(args, col, project)
	case utils.SQLServer:

		queryString = `SELECT DISTINCT C.COLUMN_NAME as 'Field', C.IS_NULLABLE as 'Null' , 
    case when C.DATA_TYPE = 'varchar' then concat(C.DATA_TYPE,'(',REPLACE(c.CHARACTER_MAXIMUM_LENGTH,'-1','max'),')') else C.DATA_TYPE end as 'Type',
    REPLACE(REPLACE(REPLACE(coalesce(C.COLUMN_DEFAULT,''),'''',''),'(',''),')','') as 'Default',
       CASE
           WHEN TC.CONSTRAINT_TYPE = 'PRIMARY KEY' THEN 'PRI'
           WHEN TC.CONSTRAINT_TYPE = 'UNIQUE' THEN 'UNI'
           ELSE isnull(TC.CONSTRAINT_TYPE,'')
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
	defer func() { _ = rows.Close() }()

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
		return result, errors.New(s.dbType + ":" + col + " not found during inspection")
	}
	return result, nil
}

func (s *SQL) getForeignKeyDetails(ctx context.Context, project, col string) ([]utils.ForeignKeysType, error) {
	queryString := ""
	args := []interface{}{}
	switch utils.DBType(s.dbType) {

	case utils.MySQL:
		queryString = "select KCU.TABLE_NAME, KCU.COLUMN_NAME, KCU.CONSTRAINT_NAME, RC.DELETE_RULE, KCU.REFERENCED_TABLE_NAME, KCU.REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KCU JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS RC ON RC.CONSTRAINT_NAME=KCU.CONSTRAINT_NAME WHERE KCU.REFERENCED_TABLE_SCHEMA = ? and KCU.TABLE_NAME = ?"
		args = append(args, project, col)
	case utils.Postgres:
		queryString = `SELECT
		tc.table_name AS "TABLE_NAME", 
		kcu.column_name AS "COLUMN_NAME", 
		tc.constraint_name AS "CONSTRAINT_NAME", 
		rc.delete_rule AS "DELETE_RULE",
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
		JOIN information_schema.referential_constraints AS rc
		  ON tc.constraint_name = rc.constraint_name
		  AND rc.unique_constraint_schema= $1
	WHERE tc.constraint_type = 'FOREIGN KEY'  AND tc.table_schema = $2  AND tc.table_name= $3
	`
		args = append(args, project, project, col)
	case utils.SQLServer:
		queryString = `SELECT 
		CCU.TABLE_NAME, CCU.COLUMN_NAME, CCU.CONSTRAINT_NAME, RC.DELETE_RULE,
		isnull(KCU.TABLE_NAME,'') AS 'REFERENCED_TABLE_NAME', isnull(KCU.COLUMN_NAME,'') AS 'REFERENCED_COLUMN_NAME'
	FROM INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE CCU
		FULL JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
			ON CCU.CONSTRAINT_NAME = RC.CONSTRAINT_NAME 
		FULL JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE KCU 
			ON KCU.CONSTRAINT_NAME = RC.CONSTRAINT_NAME  
	WHERE CCU.TABLE_SCHEMA = @p1 AND CCU.TABLE_NAME= @p2 AND KCU.TABLE_NAME= @p3`
		args = append(args, project, col, col)
	}
	rows, err := s.client.QueryxContext(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

func (s *SQL) getIndexDetails(ctx context.Context, project, col string) ([]utils.IndexType, error) {
	queryString := ""
	switch utils.DBType(s.dbType) {

	case utils.MySQL:
		queryString = `SELECT 
		TABLE_NAME, COLUMN_NAME, INDEX_NAME, SEQ_IN_INDEX, 
		(case when NON_UNIQUE = 0 then "yes" else "no" end) as IS_UNIQUE,
		(case when COLLATION = "A" then "asc" else "desc" end) as SORT 
		from INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME REGEXP '^index_'`
	case utils.Postgres:
		queryString = `select
    	t.relname as "TABLE_NAME",
    	a.attname as "COLUMN_NAME",
    	i.relname as "INDEX_NAME",
    	1 + array_position(ix.indkey, a.attnum) as "SEQ_IN_INDEX",
    	(case when ix.indisunique = false then 'no' else 'yes' end) "IS_UNIQUE",
    	(case when ix.indoption[array_position(ix.indkey, a.attnum)] = 0 then 'asc'
         when ix.indoption[array_position(ix.indkey, a.attnum)] = 3 then 'desc'
         else '' end) as "SORT"        
			from
     		pg_catalog.pg_class t
				join pg_catalog.pg_attribute a on t.oid    =      a.attrelid 
				join pg_catalog.pg_index ix    on t.oid    =     ix.indrelid
				join pg_catalog.pg_class i     on a.attnum = any(ix.indkey)
																			and i.oid    =     ix.indexrelid
				join pg_catalog.pg_namespace n on n.oid    =      t.relnamespace
				where n.nspname = $1 and t.relname = $2 and i.relname ~ '^index' and t.relkind = 'r' 
			order by
    		t.relname,
    		i.relname,
    		array_position(ix.indkey, a.attnum)`
	case utils.SQLServer:
		queryString = `SELECT 
    	TABLE_NAME = t.name,
    	COLUMN_NAME = col.name,
    	INDEX_NAME = ind.name,
    	SEQ_IN_INDEX = ic.index_column_id,
    	case when ind.is_unique = 0 then 'no' else 'yes' end as IS_UNIQUE,
    	case when ic.is_descending_key = 0 then 'asc' else 'desc' end as SORT 
			FROM 
     			sys.indexes ind 
				INNER JOIN 
     			sys.index_columns ic ON  ind.object_id = ic.object_id and ind.index_id = ic.index_id 
				INNER JOIN 
     			sys.columns col ON ic.object_id = col.object_id and ic.column_id = col.column_id 
				INNER JOIN 
     			sys.tables t ON ind.object_id = t.object_id 
				INNER JOIN 
        	sys.schemas s ON t.schema_id = s.schema_id
			WHERE 
     			ind.is_primary_key = 0  and s.name = @p1 and t.name = @p2 `
	}
	rows, err := s.client.QueryxContext(ctx, queryString, []interface{}{project, col}...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := []utils.IndexType{}
	for rows.Next() {
		indexKey := new(utils.IndexType)

		if err := rows.StructScan(indexKey); err != nil {
			return nil, err
		}

		result = append(result, *indexKey)
	}
	return result, nil
}
