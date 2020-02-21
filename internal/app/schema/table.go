package schema

import (
	"fmt"
	"sort"

	"github.com/kyleu/dbui/internal/app/conn/results"
)

type Table struct {
	Name    string
	Columns []results.Column
}

func (t *Table) AddColumn(column results.Column) {
	t.Columns = append(t.Columns, column)
}

func (t *Table) ItemID() string {
	return fmt.Sprintf("table.%s", t.Name)
}

type TableRegistry struct {
	names  []string
	tables map[string]Table
}

func (s *TableRegistry) Names() []string {
	return s.names
}

func (s *TableRegistry) Get(key string) Table {
	return s.tables[key]
}

func (s *TableRegistry) Size() int {
	return len(s.names)
}

func (s *TableRegistry) Add(t ...Table) {
	for _, x := range t {
		s.tables[x.Name] = x
	}
	var acc []string
	for _, x := range s.tables {
		acc = append(acc, x.Name)
	}
	sort.Strings(acc)
	s.names = acc
}

func newTableRegistry() TableRegistry {
	return TableRegistry{
		names:  []string{},
		tables: map[string]Table{},
	}
}
