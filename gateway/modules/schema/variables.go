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

	TableProperties struct {
		From, To     string
		Table, Field string
		DBType       string
		Group, Sort  string
		Order        int
	}
)

const (
	typeInteger        string = "Integer"
	typeString         string = "String"
	typeFloat          string = "Float"
	typeBoolean        string = "Boolean"
	typeDateTime       string = "DateTime"
	TypeID             string = "ID"
	sqlTypeIDSize      string = "50"
	typeObject         string = "Object"
	typeEnum           string = "Enum"
	directiveUnique    string = "unique"
	directiveIndex     string = "index"
	directiveForeign   string = "foreign"
	directivePrimary   string = "primary"
	directiveCreatedAt string = "createdAt"
	directiveUpdatedAt string = "updatedAt"
	directiveLink      string = "link"
	directiveDefault   string = "default"

	defaultIndexName  string = ""
	defaultIndexSort  string = "asc"
	deafultIndexOrder int    = 1
)

type indexStore []*SchemaFieldType

func (a indexStore) Len() int           { return len(a) }
func (a indexStore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a indexStore) Less(i, j int) bool { return a[i].IndexInfo.Order < a[j].IndexInfo.Order }
