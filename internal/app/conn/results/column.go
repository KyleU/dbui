package results

import "fmt"

type Column struct {
	Name     string
	T        FieldType
	Nullable bool
}

func (c Column) String() string {
	postfix := ""
	if c.Nullable {
		postfix = "+"
	}
	return fmt.Sprintf("%s:%s%s", c.Name, c.T, postfix)
}
