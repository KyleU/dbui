package results

import (
	"fmt"
	"strconv"
	"strings"
)

type Column struct {
	Name       string
	T          FieldType
	Nullable   bool
	PrimaryKey bool
	Indexed    bool
	Default    string
	Precision  int64
	Scale      int64
	Length     int64
}

func (c Column) String() string {
	return fmt.Sprintf("%s:%s", c.TypeString())
}

func (c Column) TypeString() string {
	s := c.T.Desc().Title
	if c.Length > 0 {
		s = fmt.Sprintf("%s(%v)", s, c.Length)
	}
	return s
}

func ParseArgs(t FieldType, args string) (int64, int64, int64) {
	if len(args) == 0 {
		return 0, 0, 0
	} else if strings.Contains(args, ",") {
		// TODO precision/scale
		return 0, 0, 0
	} else {
		l, err := strconv.Atoi(args)
		if err == nil {
			return 0, 0, int64(l)
		}
	}
	return 0, 0, 0
}
