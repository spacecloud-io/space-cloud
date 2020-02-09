package schema

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
		nestedObject Fields

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
	typeInteger  string = "Integer"
	typeString   string = "String"
	typeFloat    string = "Float"
	typeBoolean  string = "Boolean"
	typeDateTime string = "DateTime"
	// TypeID represents ID of schema
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

	// defaultIndexName  string = ""
	defaultIndexSort  string = "asc"
	deafultIndexOrder int    = 1
)

type indexStore []*FieldType

func (a indexStore) Len() int           { return len(a) }
func (a indexStore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a indexStore) Less(i, j int) bool { return a[i].IndexInfo.Order < a[j].IndexInfo.Order }
