package models

import "sort"

type View struct {
	Name string
}

type ViewRegistry struct {
	names []string
	views map[string]View
}

func (s *ViewRegistry) Names() []string {
	return s.names
}

func (s *ViewRegistry) Size() int {
	return len(s.names)
}

func (s *ViewRegistry) Add(t ...View) {
	for _, x := range t {
		s.views[x.Name] = x
	}
	var acc []string
	for _, x := range s.views {
		acc = append(acc, x.Name)
	}
	sort.Strings(acc)
	s.names = acc
}

func newViewRegistry() ViewRegistry {
	return ViewRegistry{
		names: []string{},
		views: map[string]View{},
	}
}
