package project

import (
	"sort"
)

type Project struct {
	Key   string
	Title string
}

type Registry struct {
	names    []string
	projects map[string]Project
}

func (s *Registry) Names() []string {
	return s.names
}

func (s *Registry) Get(key string) Project {
	return s.projects[key]
}

func (s *Registry) Size() int {
	return len(s.names)
}

func (s *Registry) Add(t ...Project) {
	for _, x := range t {
		s.projects[x.Key] = x
	}
	var acc []string
	for _, x := range s.projects {
		acc = append(acc, x.Key)
	}
	sort.Strings(acc)
	s.names = acc
}

