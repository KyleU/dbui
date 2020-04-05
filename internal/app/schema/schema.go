package schema

import "github.com/kyleu/dbui/internal/app/conn"

type Schema struct {
	ID     string
	Name   string
	Engine conn.Engine
	Tables TableRegistry
}

func NewSchema(id string, engine conn.Engine, name string) Schema {
	return Schema{
		ID:     id,
		Name:   name,
		Engine: engine,
		Tables: newTableRegistry(),
	}
}
