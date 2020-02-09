package results

type FieldType struct {
	Name string
}

var FieldTypes = map[string]FieldType{
	"string": { "string" },
	"int": { "int" },
}
