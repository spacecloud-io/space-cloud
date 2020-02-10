package model

type (

	// Type is the data structure for storing the parsed values of schema string
	Type map[string]Collection // key is database name
	// Collection is a data structure for storing fields of schema
	Collection map[string]Fields // key is collection name
	// Fields is a data structure for storing the type of field
	Fields map[string]*FieldType // key is field name
	// directiveArgs    map[string]string           // key is Directive's argument name
	// fieldType int

	// FieldType stores information about a particular column in table
	FieldType struct {
		FieldName           string
		IsFieldTypeRequired bool
		IsList              bool
		Kind                string
		//Directive           string
		NestedObject Fields

		// For directives
		IsPrimary   bool
		IsIndex     bool
		IsUnique    bool
		IsCreatedAt bool
		IsUpdatedAt bool
		IsLinked    bool
		IsForeign   bool
		IsDefault   bool
		IndexInfo   *TableProperties
		LinkedTable *TableProperties
		JointTable  *TableProperties
		Default     interface{}
	}

	// TableProperties are properties of the table
	TableProperties struct {
		From, To     string
		Table, Field string
		DBType       string
		Group, Sort  string
		Order        int
	}
)

const (
	TypeInteger        string = "Integer"
	TypeString         string = "String"
	TypeFloat          string = "Float"
	TypeBoolean        string = "Boolean"
	TypeDateTime       string = "DateTime"
	TypeID             string = "ID"
	SqlTypeIDSize      string = "50"
	TypeObject         string = "Object"
	TypeEnum           string = "Enum"
	DirectiveUnique    string = "unique"
	DirectiveIndex     string = "index"
	DirectiveForeign   string = "foreign"
	DirectivePrimary   string = "primary"
	DirectiveCreatedAt string = "createdAt"
	DirectiveUpdatedAt string = "updatedAt"
	DirectiveLink      string = "link"
	DirectiveDefault   string = "default"

	//DefaultIndexName  string = ""
	DefaultIndexSort  string = "asc"
	DefaultIndexOrder int    = 1
)
