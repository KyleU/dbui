package models

import "sort"

type Table struct {
	Name string
}

type TableRegistry struct {
	names  []string
	tables map[string]Table
}

func (s *TableRegistry) Names() []string {
	return s.names
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
