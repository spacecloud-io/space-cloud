package schema

type (

	// schemaType is the data structure for storing the parsed values of schema string
	schemaType       map[string]schemaCollection // key is database name
	schemaCollection map[string]SchemaFields     // key is collection name
	SchemaFields     map[string]*SchemaFieldType // key is field name
	directiveArgs    map[string]string           // key is Directive's argument name
	fieldType        int

	SchemaFieldType struct {
		FieldName           string
		IsFieldTypeRequired bool
		IsList              bool
		Kind                string
		//Directive           string
		nestedObject SchemaFields

		// For directives
		IsPrimary   bool
		IsUnique    bool
		IsCreatedAt bool
		IsUpdatedAt bool
		IsLinked    bool
		IsForeign   bool
		LinkedTable *TableProperties
		JointTable  *TableProperties
	}
	TableProperties struct {
		From, To     string
		Table, Field string
		DBType       string
	}
)

const (
	typeInteger   string = "Integer"
	typeString    string = "String"
	typeFloat     string = "Float"
	typeBoolean   string = "Boolean"
	typeDateTime  string = "DateTime"
	TypeID        string = "ID"
	sqlTypeIDSize string = "50"
	typeObject    string = "Object"

	directiveUnique    string = "unique"
	directiveForeign   string = "foreign"
	directivePrimary   string = "primary"
	directiveCreatedAt string = "createdAt"
	directiveUpdatedAt string = "updatedAt"
	directiveLink      string = "link"
)
