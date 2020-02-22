package util

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"regexp"
)

func MicrosToMillis(l language.Tag, i int) string {
	ms := i / 1000
	if ms >= 20 {
		return FormatInteger(l, ms) + "ms"
	}
	x := float64(ms) + (float64(i%1000) / 1000)
	p := message.NewPrinter(l)
	return p.Sprintf("%.3f", x) + "ms"
}

func FormatInteger(l language.Tag, v int) string {
	p := message.NewPrinter(l)
	return p.Sprintf("%d", v)
}

func PluralChoice(plural string, single string, v int) string {
	if v == 1 || v == -1 {
		return single
	}
	return plural
}

var re *regexp.Regexp

func PathParams(s string) []string {
	if re == nil {
		re, _ = regexp.Compile("{([^}]*)}")
	}
	matches := re.FindAll([]byte(s), -1)

	ret := make([]string, 0)
	for _, m := range matches {
		ret = append(ret, string(m))
	}

	return ret
}
