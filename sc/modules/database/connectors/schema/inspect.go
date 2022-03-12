package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spaceuptech/helpers"
)

// ParseCollectionDescription coverts the output of descibe table to sc schema description format
func ParseCollectionDescription(dbType, col string, fields []model.InspectorFieldType, indexes []model.IndexType, existingSchema model.CollectionSchemas) (model.CollectionSchemas, error) {
	currentSchema := model.CollectionSchemas{}
	currentFieldSchemas := model.FieldSchemas{}

	for _, field := range fields {
		fieldSchema := model.FieldType{FieldName: field.ColumnName}

		// Check if field is required?
		if field.FieldNull == "NO" {
			fieldSchema.IsFieldTypeRequired = true
		}

		// Then we inpect the fieldtype based on the database type
		switch model.DBType(dbType) {
		case model.Postgres:
			if err := inspectionPostgresCheckFieldType(col, field, &fieldSchema); err != nil {
				return nil, err
			}
		case model.MySQL:
			if err := inspectionMySQLCheckFieldType(col, field, &fieldSchema); err != nil {
				return nil, err
			}
		case model.SQLServer:
			if err := inspectionSQLServerCheckFieldType(col, field, &fieldSchema); err != nil {
				return nil, err
			}
		}

		// Check if the field has a default value present
		if field.FieldDefault != "" {
			fieldSchema.IsDefault = true

			// Process the default values to a compatible format
			switch fieldSchema.Kind {
			case model.TypeString, model.TypeVarChar, model.TypeChar, model.TypeID, model.TypeDateTime, model.TypeDate, model.TypeTime, model.TypeDateTimeWithZone:
				field.FieldDefault = fmt.Sprintf("\"%s\"", field.FieldDefault)
			case model.TypeJSON:
				data, err := json.Marshal(field.FieldDefault)
				if err != nil {
					return nil, helpers.Logger.LogError("generate-inspection", "Unable to parse column having a default value of type json", err, nil)
				}
				field.FieldDefault = string(data)
			}
			fieldSchema.Default = field.FieldDefault
		}

		// Check if foreign key exists
		if field.RefTableName != "" && field.RefColumnName != "" {
			fieldSchema.IsForeign = true
			fieldSchema.JointTable = &model.TableProperties{Table: field.RefTableName, To: field.RefColumnName, OnDelete: field.DeleteRule, ConstraintName: field.ConstraintName}
		}

		// Check if field is of type autoincrement
		if field.AutoIncrement == "true" {
			fieldSchema.IsAutoIncrement = true
			if model.DBType(dbType) == model.Postgres && strings.HasPrefix(field.FieldDefault, "nextval") {
				// override the default value, this is a special case if a postgres column has a auto increment value, the default value that database returns is -> ( nextval(auto_increment_test_auto_increment_test_seq )
				fieldSchema.Default = ""
				fieldSchema.IsDefault = false
			}
		}

		// Check if the field has any indexes on it
		for _, indexValue := range indexes {
			if indexValue.ColumnName == field.ColumnName {
				temp := &model.TableProperties{Order: indexValue.Order, Sort: indexValue.Sort, ConstraintName: indexValue.IndexName}
				if indexValue.IsPrimary {
					fieldSchema.IsPrimary = true
					fieldSchema.PrimaryKeyInfo = &model.TableProperties{
						Order: indexValue.Order,
					}
					continue
				} else if indexValue.IsUnique {
					temp.IsUnique = true
				} else {
					temp.IsIndex = true
				}

				if fieldSchema.IndexInfo == nil {
					fieldSchema.IndexInfo = make([]*model.TableProperties, 0)
				}
				if strings.HasPrefix(indexValue.IndexName, "index__") {
					indexValue.IndexName = getGroupNameFromIndexName(indexValue.IndexName)
				}
				temp.Group = indexValue.IndexName
				temp.Field = indexValue.ColumnName
				// index is created through gateway, as it follows our naming convention
				fieldSchema.IndexInfo = append(fieldSchema.IndexInfo, temp)
			}
		}
		// field name
		currentFieldSchemas[field.ColumnName] = &fieldSchema
	}

	currentSchema[col] = currentFieldSchemas

	// Some attributes of the schema can only be filled from the existing schema.
	// Check if the existing schmea
	for columnName, existingFieldSchema := range existingSchema[col] {
		// Add any previous links that were present
		if existingFieldSchema.IsLinked {
			currentFieldSchemas[columnName] = existingFieldSchema
			continue
		}

		// Check if the current schema contains the column defined in the previous schema
		currentTableInfo, ok := currentFieldSchemas[columnName]
		if !ok {
			continue
		}

		// Basic updates
		if existingFieldSchema.Kind == model.TypeID {
			currentTableInfo.Kind = model.TypeID
		}
		if existingFieldSchema.IsCreatedAt {
			currentTableInfo.IsCreatedAt = true
		}
		if existingFieldSchema.IsUpdatedAt {
			currentTableInfo.IsUpdatedAt = true
		}
	}

	return currentSchema, nil
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
