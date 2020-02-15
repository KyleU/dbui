package models

type Schema struct {
	Id     string
	Name   string
	Tables TableRegistry
	Views  ViewRegistry
}

func NewSchema(id string, name string) Schema {
	return Schema{
		Id:     id,
		Name:   name,
		Tables: newTableRegistry(),
		Views:  newViewRegistry(),
	}
}
