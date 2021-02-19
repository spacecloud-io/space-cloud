package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	schemaHelpers "github.com/spaceuptech/space-cloud/gateway/modules/schema/helpers"
)

// SchemaInspection returns the schema in schema definition language (SDL)
func (s *Schema) SchemaInspection(ctx context.Context, dbAlias, project, col string, realSchema model.Collection) (string, error) {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return "", err
	}

	if dbType == "mongo" {
		return "", nil
	}

	inspectionCollection, err := s.Inspector(ctx, dbAlias, dbType, project, col, realSchema)
	if err != nil {
		return "", err
	}

	return generateSDL(inspectionCollection)

}

// Inspector generates schema
func (s *Schema) Inspector(ctx context.Context, dbAlias, dbType, project, col string, realSchema model.Collection) (model.Collection, error) {
	fields, indexes, err := s.crud.DescribeTable(ctx, dbAlias, col)

	if err != nil {
		return nil, err
	}
	currentSchema, err := generateInspection(dbType, col, fields, indexes)
	if err != nil {
		return nil, err
	}
	currentTableFields, ok := currentSchema[col]
	if !ok {
		return currentSchema, nil
	}

	for columnName, realColumnInfo := range realSchema[col] {
		if realColumnInfo.IsLinked {
			currentTableFields[columnName] = realColumnInfo
			continue
		}

		currentTableInfo, ok := currentTableFields[columnName]
		if !ok {
			continue
		}
		if realColumnInfo.Kind == model.TypeID {
			currentTableInfo.Kind = model.TypeID
		}

		if realColumnInfo.IsCreatedAt {
			currentTableInfo.IsCreatedAt = true
		}
		if realColumnInfo.IsUpdatedAt {
			currentTableInfo.IsUpdatedAt = true
		}
	}

	return currentSchema, nil
}

func generateInspection(dbType, col string, fields []model.InspectorFieldType, indexes []model.IndexType) (model.Collection, error) {
	inspectionCollection := model.Collection{}
	inspectionFields := model.Fields{}

	for _, field := range fields {
		fieldDetails := model.FieldType{FieldName: field.ColumnName}

		// check if field nullable (!)
		if field.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}

		// field type
		switch model.DBType(dbType) {
		case model.Postgres:
			if err := inspectionPostgresCheckFieldType(col, field, &fieldDetails); err != nil {
				return nil, err
			}
		case model.MySQL:
			if err := inspectionMySQLCheckFieldType(col, field, &fieldDetails); err != nil {
				return nil, err
			}
		case model.SQLServer:
			if err := inspectionSQLServerCheckFieldType(col, field, &fieldDetails); err != nil {
				return nil, err
			}
		}

		// default key
		if field.FieldDefault != "" {
			fieldDetails.IsDefault = true

			// add string between quotes
			switch fieldDetails.Kind {
			case model.TypeString, model.TypeVarChar, model.TypeChar, model.TypeID, model.TypeDateTime, model.TypeDate, model.TypeTime, model.TypeDateTimeWithZone:
				field.FieldDefault = fmt.Sprintf("\"%s\"", field.FieldDefault)
			case model.TypeJSON:
				data, err := json.Marshal(field.FieldDefault)
				if err != nil {
					return nil, helpers.Logger.LogError("generate-inspection", "Unable to parse column having a default value of type json", err, nil)
				}
				field.FieldDefault = string(data)
			}
			fieldDetails.Default = field.FieldDefault
		}

		// check foreignKey & identify if relation exists
		if field.RefTableName != "" && field.RefColumnName != "" {
			fieldDetails.IsForeign = true
			fieldDetails.JointTable = &model.TableProperties{Table: field.RefTableName, To: field.RefColumnName, OnDelete: field.DeleteRule, ConstraintName: field.ConstraintName}
		}

		if field.AutoIncrement == "true" {
			fieldDetails.IsAutoIncrement = true
			if model.DBType(dbType) == model.Postgres && strings.HasPrefix(field.FieldDefault, "nextval") {
				// override the default value, this is a special case if a postgres column has a auto increment value, the default value that database returns is -> ( nextval(auto_increment_test_auto_increment_test_seq )
				fieldDetails.Default = ""
				fieldDetails.IsDefault = false
			}
		}

		for _, indexValue := range indexes {
			if indexValue.ColumnName == field.ColumnName {
				temp := &model.TableProperties{Order: indexValue.Order, Sort: indexValue.Sort, ConstraintName: indexValue.IndexName}
				if indexValue.IsPrimary {
					fieldDetails.IsPrimary = true
					fieldDetails.PrimaryKeyInfo = &model.TableProperties{
						Order: indexValue.Order,
					}
					continue
				} else if indexValue.IsUnique {
					temp.IsUnique = true
				} else {
					temp.IsIndex = true
				}

				if fieldDetails.IndexInfo == nil {
					fieldDetails.IndexInfo = make([]*model.TableProperties, 0)
				}
				if strings.HasPrefix(indexValue.IndexName, "index__") {
					indexValue.IndexName = getGroupNameFromIndexName(indexValue.IndexName)
				}
				temp.Group = indexValue.IndexName
				temp.Field = indexValue.ColumnName
				// index is created through gateway, as it follows our naming convention
				fieldDetails.IndexInfo = append(fieldDetails.IndexInfo, temp)
			}
		}
		// field name
		inspectionFields[field.ColumnName] = &fieldDetails
	}

	if len(inspectionFields) != 0 {
		inspectionCollection[col] = inspectionFields
	}
	return inspectionCollection, nil
}

func getGroupNameFromIndexName(indexName string) string {
	// ignoring the length check as the length is assured to be 3
	return strings.Split(indexName, "__")[2]
}

func inspectionMySQLCheckFieldType(col string, field model.InspectorFieldType, fieldDetails *model.FieldType) error {
	result := strings.Split(field.FieldType, "(")

	switch result[0] {
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "varchar":
		fieldDetails.Kind = model.TypeVarChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "char":
		fieldDetails.Kind = model.TypeChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "tinytext", "text", "mediumtext", "longtext":
		fieldDetails.Kind = model.TypeString
	case "smallint":
		fieldDetails.Kind = model.TypeSmallInteger
	case "bigint":
		fieldDetails.Kind = model.TypeBigInteger
	case "mediumint", "int":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double":
		fieldDetails.Kind = model.TypeFloat
	case "decimal":
		fieldDetails.Kind = model.TypeDecimal
		if field.NumericPrecision > 0 || field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.NumericPrecision,
				Scale:     field.NumericScale,
			}
		}
	case "datetime":
		fieldDetails.Kind = model.TypeDateTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "timestamp":
		fieldDetails.Kind = model.TypeDateTimeWithZone
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "bit", "tinyint":
		fieldDetails.Kind = model.TypeBoolean
	case "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return helpers.Logger.LogError("", fmt.Sprintf("Cannot track/inspect table (%s)", col), fmt.Errorf("table contains a column (%s) with type (%s) which is not supported by space cloud", fieldDetails.FieldName, result), nil)
	}
	return nil
}

func inspectionSQLServerCheckFieldType(col string, field model.InspectorFieldType, fieldDetails *model.FieldType) error {
	result := strings.Split(field.FieldType, "(")

	switch result[0] {
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "varchar", "nvarchar":
		if field.VarcharSize == -1 {
			fieldDetails.Kind = model.TypeString
			return nil
		}
		fieldDetails.Kind = model.TypeVarChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "char", "nchar":
		fieldDetails.Kind = model.TypeChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "text", "ntext":
		fieldDetails.Kind = model.TypeString
	case "smallint":
		fieldDetails.Kind = model.TypeSmallInteger
	case "bigint":
		fieldDetails.Kind = model.TypeBigInteger
	case "int":
		fieldDetails.Kind = model.TypeInteger
	case "numeric", "decimal":
		fieldDetails.Kind = model.TypeDecimal
		if field.NumericPrecision > 0 || field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.NumericPrecision,
				Scale:     field.NumericScale,
			}
		}
	case "float", "real":
		fieldDetails.Kind = model.TypeFloat
	case "datetime", "datetime2", "smalldatetime":
		fieldDetails.Kind = model.TypeDateTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "datetimeoffset":
		fieldDetails.Kind = model.TypeDateTimeWithZone
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "bit", "tinyint":
		fieldDetails.Kind = model.TypeBoolean
	case "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return helpers.Logger.LogError("", fmt.Sprintf("Cannot track/inspect table (%s)", col), fmt.Errorf("table contains a column (%s) with type (%s) which is not supported by space cloud", fieldDetails.FieldName, result), nil)
	}
	return nil
}

func inspectionPostgresCheckFieldType(col string, field model.InspectorFieldType, fieldDetails *model.FieldType) error {
	result := strings.Split(field.FieldType, "(")

	switch result[0] {
	case "uuid":
		fieldDetails.Kind = model.TypeUUID
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time without time zone", "time with time zone":
		fieldDetails.Kind = model.TypeTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "character varying":
		fieldDetails.Kind = model.TypeVarChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "character":
		fieldDetails.Kind = model.TypeChar
		fieldDetails.TypeIDSize = field.VarcharSize
	case "text", "name":
		fieldDetails.Kind = model.TypeString
	case "integer", "serial":
		fieldDetails.Kind = model.TypeInteger
	case "smallint", "smallserial":
		fieldDetails.Kind = model.TypeSmallInteger
	case "bigint", "bigserial":
		fieldDetails.Kind = model.TypeBigInteger
	case "numeric":
		fieldDetails.Kind = model.TypeDecimal
		if field.NumericPrecision > 0 || field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.NumericPrecision,
				Scale:     field.NumericScale,
			}
		}
	case "real", "double precision":
		fieldDetails.Kind = model.TypeFloat
	case "timestamp without time zone":
		fieldDetails.Kind = model.TypeDateTime
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "timestamp with time zone":
		fieldDetails.Kind = model.TypeDateTimeWithZone
		if field.DateTimePrecision > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Precision: field.DateTimePrecision,
			}
		}
	case "boolean":
		fieldDetails.Kind = model.TypeBoolean
	case "jsonb", "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return helpers.Logger.LogError("", fmt.Sprintf("Cannot track/inspect table (%s)", col), fmt.Errorf("table contains a column (%s) with type (%s) which is not supported by space cloud", fieldDetails.FieldName, result), nil)
	}
	return nil
}

// GetCollectionSchema returns schemas of collection aka tables for specified project & database
func (s *Schema) GetCollectionSchema(ctx context.Context, project, dbAlias string) (map[string]*config.TableRule, error) {

	collections := []string{}
	for _, dbSchema := range s.dbSchemas {
		if dbSchema.DbAlias == dbAlias {
			collections = append(collections, dbSchema.Table)
			break
		}
	}

	parsedSchema, _ := schemaHelpers.Parser(s.dbSchemas)
	projectConfig := config.Crud{}
	projectConfig[dbAlias] = &config.CrudStub{}
	for _, colName := range collections {
		if colName == "default" {
			continue
		}
		schema, err := s.SchemaInspection(ctx, dbAlias, project, colName, parsedSchema[dbAlias])
		if err != nil {
			return nil, err
		}

		if projectConfig[dbAlias].Collections == nil {
			projectConfig[dbAlias].Collections = map[string]*config.TableRule{}
		}
		projectConfig[dbAlias].Collections[colName] = &config.TableRule{Schema: schema}
	}
	return projectConfig[dbAlias].Collections, nil
}
