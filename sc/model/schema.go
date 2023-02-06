package model

type (

	// DBSchemas is the data structure for storing the parsed values of schema string
	DBSchemas map[string]CollectionSchemas // key is database name
	// CollectionSchemas is a data structure for storing fields of schema
	CollectionSchemas map[string]FieldSchemas // key is collection name
	// FieldSchemas is a data structure for storing the type of field
	FieldSchemas map[string]*FieldType // key is field name

	// FieldType stores information about a particular column in table
	FieldType struct {
		FieldName           string     `json:"fieldName"`
		IsFieldTypeRequired bool       `json:"isFieldTypeRequired"`
		IsList              bool       `json:"isList"`
		Kind                string     `json:"kind"`
		Args                *FieldArgs `json:"args"`
		// Directive           string
		NestedObject FieldSchemas `json:"nestedObject"`
		IsPrimary    bool         `json:"isPrimary"`
		// For directives
		IsCreatedAt     bool `json:"isCreatedAt"`
		IsUpdatedAt     bool `json:"isUpdatedAt"`
		IsLinked        bool `json:"isLinked"`
		IsForeign       bool `json:"isForeign"`
		IsDefault       bool `json:"isDefault"`
		IsAutoIncrement bool
		PrimaryKeyInfo  *TableProperties   `json:"primaryKeyInfo"`
		IndexInfo       []*TableProperties `json:"indexInfo"`
		LinkedTable     *LinkProperties    `json:"linkedTable"`
		JointTable      *TableProperties   `json:"jointTable"`
		Default         interface{}        `json:"default"`
		TypeIDSize      int                `json:"size"`
	}

	// FieldArgs are properties of the column
	FieldArgs struct {
		// Precision is used to hold precision information for data types Float
		// It represent number the digits to be stored
		Precision int `json:"precision"`
		// Scale is used to hold scale information for data types Float,Time,DateTime
		// It represent number the digits to be stored after decimal
		Scale int `json:"scale"`
	}

	// LinkProperties describes the properties of a link
	LinkProperties struct {
		DB       string
		Table    string
		Field    string // Used for modeling many to many relationships
		From, To string
	}

	// TableProperties are properties of the table
	TableProperties struct {
		// IsIndex tells us if this is an indexed column
		IsIndex bool `json:"isIndex"`
		// IsUnique tells us if this is an unique indexed column
		IsUnique       bool `json:"isUnique"`
		From           string
		To             string
		On             map[string]interface{}
		Table          string
		Field          string
		OnDelete       string
		DBType         string
		Group          string
		Sort           string
		Order          int
		ConstraintName string
	}
)

const (
	// TypeDate is variable used for Variable of type Date
	TypeDate string = "Date"
	// TypeTime is variable used for Variable of type Time
	TypeTime string = "Time"
	// TypeUUID is variable used for Variable of type UUID
	TypeUUID string = "UUID"
	// TypeInteger is variable used for Variable of type Integer
	TypeInteger string = "Integer"
	// TypeSmallInteger is variable used for Variable of type small int
	TypeSmallInteger string = "SmallInteger"
	// TypeBigInteger is variable used for Variable of type big int
	TypeBigInteger string = "BigInteger"
	// TypeChar is variable used for Variable of type characters with fixed size
	TypeChar string = "Char"
	// TypeVarChar is variable used for Variable of type characters with variable size
	TypeVarChar string = "Varchar"
	// TypeString is variable used for Variable of type String
	TypeString string = "String"
	// TypeFloat is a data type used for storing fractional values without specifying the precision,
	// the precision is set by the database
	TypeFloat string = "Float"
	// TypeDecimal is a data type used for storing fractional values in which the precision can be specified by the user
	TypeDecimal string = "Decimal"
	// TypeBoolean is variable used for Variable of type Boolean
	TypeBoolean string = "Boolean"
	// TypeDateTimeWithZone is variable used for Variable of type DateTime
	TypeDateTimeWithZone string = "DateTimeWithZone"
	// TypeDateTime is variable used for Variable of type DateTime
	TypeDateTime string = "DateTime"
	// TypeID is variable used for Variable of type ID
	TypeID string = "ID"
	// TypeJSON is variable used for Variable of type Jsonb
	TypeJSON string = "JSON"
	// DefaultCharacterSize is variable used for specifying size of sql type ID
	DefaultCharacterSize int = 100
	// TypeObject is a string with value object
	TypeObject string = "Object"
	// TypeEnum is a variable type enum
	TypeEnum string = "Enum"
	// DirectiveUnique is used in schema module to add unique index
	DirectiveUnique string = "unique"
	// DirectiveIndex is used in schema module to add index
	DirectiveIndex string = "index"
	// DirectiveForeign is used in schema module to add foreign key
	DirectiveForeign string = "foreign"
	// DirectivePrimary is used in schema module to add primary key
	DirectivePrimary string = "primary"
	// DirectiveAutoIncrement is used in schema module to add primary key
	DirectiveAutoIncrement string = "autoIncrement"
	// DirectiveCreatedAt is used in schema module to specify the created location
	DirectiveCreatedAt string = "createdAt"
	// DirectiveUpdatedAt  is used in schema module to add Updated location
	DirectiveUpdatedAt string = "updatedAt"
	// DirectiveLink is used in schema module to add link
	DirectiveLink string = "link"
	// DirectiveDefault is used to add default key
	DirectiveDefault string = "default"
	// DirectiveArgs is used in schema module to specify the created location
	DirectiveArgs string = "args"
	// DirectiveStringSize denotes the maximum allowable character for field type Char, Varchar, ID
	DirectiveStringSize string = "size"

	// DefaultIndexSort specifies default order of sorting
	DefaultIndexSort string = "asc"
	// DefaultIndexOrder specifies default order of order
	DefaultIndexOrder int = 1

	// DefaultScale specifies the default scale to be used for sql column types float,date,datetime if not provided
	DefaultScale int = 10
	// DefaultPrecision specifies the default precision to be used for sql column types float if not provided
	DefaultPrecision int = 38
	// DefaultDateTimePrecision specifies the default precision to be used for sql column types datetime & time
	DefaultDateTimePrecision int = 6
)

// InspectorFieldType is the type for storing sql inspection information
type InspectorFieldType struct {
	// TableSchema is the schema name for postgres & sqlserver.
	// it is the database name for mysql
	TableSchema string `db:"TABLE_SCHEMA"`
	TableName   string `db:"TABLE_NAME"`

	ColumnName string `db:"COLUMN_NAME"`
	// FieldType is the data type of column
	FieldType string `db:"DATA_TYPE"`
	// FieldNull specifies whether the given column can be null or not
	// It can be either (NO) or (YES)
	FieldNull       string `db:"IS_NULLABLE"`
	OrdinalPosition string `db:"ORDINAL_POSITION"`
	// FieldDefault specifies the default value of columns
	FieldDefault string `db:"DEFAULT"`
	// AutoIncrement specifies whether the column has auto increment constraint
	// It can be either (true) or (false)
	AutoIncrement     string `db:"AUTO_INCREMENT"`
	VarcharSize       int    `db:"CHARACTER_MAXIMUM_LENGTH"`
	NumericScale      int    `db:"NUMERIC_SCALE"`
	NumericPrecision  int    `db:"NUMERIC_PRECISION"`
	DateTimePrecision int    `db:"DATETIME_PRECISION"`

	ConstraintName string `db:"CONSTRAINT_NAME"`
	DeleteRule     string `db:"DELETE_RULE"`
	RefTableSchema string `db:"REFERENCED_TABLE_SCHEMA"`
	RefTableName   string `db:"REFERENCED_TABLE_NAME"`
	RefColumnName  string `db:"REFERENCED_COLUMN_NAME"`
}

// IndexType is the type use to indexkey information of sql inspection
type IndexType struct {
	TableSchema string `db:"TABLE_SCHEMA"`
	TableName   string `db:"TABLE_NAME"`
	ColumnName  string `db:"COLUMN_NAME"`
	IndexName   string `db:"INDEX_NAME"`
	Order       int    `db:"SEQ_IN_INDEX"`
	// Sort can be either (asc) or (desc)
	Sort string `db:"SORT"`
	// IsUnique specifies whether the column has a unique index
	IsUnique bool `db:"IS_UNIQUE"`
	// IsPrimary specifies whether the column has a index
	IsPrimary bool `db:"IS_PRIMARY"`
}
