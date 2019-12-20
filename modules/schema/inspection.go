package schema

import (
	"context"
	"errors"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// SchemaInspection returns the schema in schema definition language (SDL)
func (s *Schema) SchemaInspection(ctx context.Context, dbAlias, project, col string) (string, error) {
	if dbAlias == "mongo" {
		return "", nil
	}

	inspectionCollection, err := s.Inspector(ctx, dbAlias, project, col)
	if err != nil {
		return "", err
	}

	return generateSDL(inspectionCollection)

}

// Inspector does something
func (s *Schema) Inspector(ctx context.Context, dbAlias, project, col string) (schemaCollection, error) {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return nil, err
	}
	fields, foreignkeys, err := s.crud.DescribeTable(ctx, dbType, project, col)
	if err != nil {
		return nil, err
	}
	inspectionCollection := schemaCollection{}
	inspectionFields := SchemaFields{}

	for _, field := range fields {
		fieldDetails := SchemaFieldType{FieldName: field.FieldName}

		// check if field nullable (!)
		if field.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}

		// field type
		if utils.DBType(dbType) == utils.Postgres {
			if err := inspectionPostgresCheckFieldType(field.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		} else {
			if err := inspectionMySQLCheckFieldType(field.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		}
		// check if list
		if field.FieldKey == "PRI" {
			fieldDetails.IsPrimary = true
		}

		if field.FieldKey == "UNI" {
			fieldDetails.IsUnique = true
		}

		// check foreignKey & identify if relation exists
		for _, foreignValue := range foreignkeys {
			if foreignValue.ColumnName == field.FieldName && foreignValue.RefTableName != "" && foreignValue.RefColumnName != "" {
				fieldDetails.IsForeign = true
				fieldDetails.JointTable = &TableProperties{Table: foreignValue.RefTableName, To: foreignValue.RefColumnName}
			}
		}

		// field name
		inspectionFields[field.FieldName] = &fieldDetails
	}
	if len(inspectionFields) != 0 {
		inspectionCollection[col] = inspectionFields
	}
	return inspectionCollection, nil
}

func inspectionMySQLCheckFieldType(typeName string, fieldDetails *SchemaFieldType) error {
	if typeName == "varchar("+sqlTypeIDSize+")" {
		fieldDetails.Kind = TypeID
		return nil
	}

	result := strings.Split(typeName, "(")

	switch result[0] {
	case "varchar":
		fieldDetails.Kind = TypeID // for sql server
	case "char", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal":
		fieldDetails.Kind = typeString
	case "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = typeInteger
	case "float", "double":
		fieldDetails.Kind = typeFloat
	case "date", "time", "datetime", "timestamp":
		fieldDetails.Kind = typeDateTime
	case "tinyint", "boolean":
		fieldDetails.Kind = typeBoolean
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

func inspectionPostgresCheckFieldType(typeName string, fieldDetails *SchemaFieldType) error {
	if typeName == "character varying("+sqlTypeIDSize+")" {
		fieldDetails.Kind = TypeID
		return nil
	}

	result := strings.Split(typeName, " ")
	result = strings.Split(result[0], "(")

	switch result[0] {
	case "character", "bit", "text":
		fieldDetails.Kind = typeString
	case "bigint", "bigserial", "integer", "numeric", "smallint", "smallserial", "serial":
		fieldDetails.Kind = typeInteger
	case "float", "double", "real":
		fieldDetails.Kind = typeFloat
	case "date", "time", "datetime", "timestamp", "interval":
		fieldDetails.Kind = typeDateTime
	case "boolean":
		fieldDetails.Kind = typeBoolean

	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

// GetCollectionSchema returns schemas of collection aka tables for specified project & database
func (s *Schema) GetCollectionSchema(ctx context.Context, project, dbType string) (map[string]*config.TableRule, error) {

	collections := []string{}
	for dbName, crudValue := range s.config {
		if dbName == dbType {
			for colName := range crudValue.Collections {
				collections = append(collections, colName)
			}
			break
		}
	}

	projectConfig := config.Crud{}
	projectConfig[dbType] = &config.CrudStub{}
	for _, colName := range collections {
		if colName == "default" {
			continue
		}
		schema, err := s.SchemaInspection(ctx, dbType, project, colName)
		if err != nil {
			return nil, err
		}

		if projectConfig[dbType].Collections == nil {
			projectConfig[dbType].Collections = map[string]*config.TableRule{}
		}
		projectConfig[dbType].Collections[colName] = &config.TableRule{Schema: schema}
	}
	return projectConfig[dbType].Collections, nil
}
