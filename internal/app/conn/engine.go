package conn

type Engine struct {
	Key  string
	Name string
}

var PostgreSQL = Engine{
	Key:  "pgx",
	Name: "PostgreSQL",
}

var MySQL = Engine{
	Key:  "mysql",
	Name: "MySQL",
}

var SQLite = Engine{
	Key:  "sqlite3",
	Name: "SQLite",
}

var AllEngines = []Engine{PostgreSQL, MySQL, SQLite}

func (t Engine) String() string {
	return t.Key
}

func EngineFromString(s string) Engine {
	for _, t := range AllEngines {
		if t.String() == s {
			return t
		}
	}
	return PostgreSQL
}
