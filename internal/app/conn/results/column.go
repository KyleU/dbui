package results

import "fmt"

type Column struct {
	Name string
	T    FieldType
}

func (c Column) String() string {
	return fmt.Sprintf("%s:%s", c.Name, c.T)
}
