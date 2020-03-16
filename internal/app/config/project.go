package config

import (
	"github.com/gofrs/uuid"
	"sort"
)

type Project struct {
	Key         string    `db:"key"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Owner       uuid.UUID `db:"owner"`
	Engine      string    `db:"engine"`
	Url         string    `db:"url"`
}

type ProjectRegistry struct {
	names    []string
	projects map[string]Project
}

func (s *ProjectRegistry) Names() []string {
	return s.names
}

func (s *ProjectRegistry) Get(key string) Project {
	return s.projects[key]
}

func (s *ProjectRegistry) Size() int {
	return len(s.names)
}

func (s *ProjectRegistry) Add(t ...Project) {
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
