package schema

type Schema struct {
	ID     string
	Name   string
	Engine string
	Tables TableRegistry
}

func NewSchema(id string, name string) Schema {
	return Schema{
		ID:     id,
		Name:   name,
		Engine: "pgx",
		Tables: newTableRegistry(),
	}
}
