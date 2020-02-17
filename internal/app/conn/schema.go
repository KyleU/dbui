package conn

import "github.com/kyleu/dbui/internal/app/util"

func LoadSchema(ai util.AppInfo, id string, url string) {
	rs, err := GetResult(ai.Logger, "", "")
	if err != nil {
		ai.Logger.Warn("Error loading schema")
	} else {
		println(rs)
	}
}
