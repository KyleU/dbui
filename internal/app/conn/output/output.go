package output

import (
	"github.com/kyleu/dbui/internal/app/conn/results"
	"strings"
)

func OutputFor(result *results.ResultSet, out string) (string, error) {
	switch getFormat(out) {
	case "string":
		return AsString(result)
	default:
		return AsTable(result)
	}
}

func getFormat(o string) string {
	switch strings.ToLower(o) {
	case "string":
		return "string"
	case "csv":
		return "csv"
	case "json":
		return "json"
	case "markdown":
		return "markdown"
	default:
		return "table"
	}
}
