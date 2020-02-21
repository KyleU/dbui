package schema

type Schema struct {
	ID     string
	Name   string
	Tables TableRegistry
	Views  ViewRegistry
}

func NewSchema(id string, name string) Schema {
	return Schema{
		ID:     id,
		Name:   name,
		Tables: newTableRegistry(),
		Views:  newViewRegistry(),
	}
}
