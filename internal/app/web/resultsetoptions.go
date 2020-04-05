package web

import (
	"strings"

	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
)

type ResultOptions struct {
	Profile  util.UserProfile
	Engine   conn.Engine
	Sortable bool
	SortCol  string
	SortAsc  bool
}

func NewResultOptions(profile util.UserProfile, engine conn.Engine, sortable bool) ResultOptions {
	return ResultOptions{Profile: profile, Engine: engine, Sortable: sortable, SortCol: "", SortAsc: true}
}

func (opts *ResultOptions) SortIconFor(name string) string {
	if opts.SortCol == name && opts.SortAsc {
		return "icon: triangle-up"
	}
	return "icon: triangle-down"
}

func (opts *ResultOptions) CurrentIconFor() string {
	if opts.SortAsc {
		return "icon: triangle-down"
	}
	return "icon: triangle-up"
}

func (opts *ResultOptions) SortTitleFor(name string) string {
	if opts.SortCol == name {
		if opts.SortAsc {
			return ", sorted ascending"
		}
		return ", sorted descending"
	}
	return ""
}

func FromQueryString(profile util.UserProfile, sortable bool, q map[string][]string) ResultOptions {
	col, ok := q["sc"]
	sortCol := ""
	sortAsc := true
	if ok && len(col) > 0 {
		sortCol = col[0]
		order, ok := q["so"]
		if ok && len(order) > 0 {
			sortAsc = order[0] == "a"
		}
	}
	return ResultOptions{
		Profile:  profile,
		Sortable: sortable,
		SortCol:  sortCol,
		SortAsc:  sortAsc,
	}
}

func (opts *ResultOptions) ToQueryString(nameOverride string) string {
	ret := make([]string, 0)

	if nameOverride != "" {
		ret = append(ret, "sc="+nameOverride)
		if opts.SortCol == nameOverride && opts.SortAsc {
			ret = append(ret, "so=d")
		} else {
			ret = append(ret, "so=a")
		}
	} else if opts.SortCol != "" {
		ret = append(ret, "sc="+opts.SortCol)
		if opts.SortAsc {
			ret = append(ret, "so=d")
		} else {
			ret = append(ret, "so=a")
		}
	}
	prefix := ""
	if len(ret) > 0 {
		prefix = "?"
	}
	return prefix + strings.Join(ret, "&")
}

func (opts *ResultOptions) ToSQL(name string) string {
	sb := &strings.Builder{}
	sb.WriteString("select * from ")
	switch opts.Engine {
	case conn.PostgreSQL:
		sb.WriteString("\"")
		sb.WriteString(name)
		sb.WriteString("\"")
	default:
		sb.WriteString(name)
	}
	if opts.SortCol != "" {
		sb.WriteString(" order by \"")
		sb.WriteString(opts.SortCol)
		sb.WriteString("\"")
		if !opts.SortAsc {
			sb.WriteString(" desc")
		}
	}
	return sb.String()
}
