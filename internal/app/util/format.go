package util

import "fmt"

func MicrosToMillis(i int64) string {
	ms := i / 1000
	if ms >= 20 {
		return fmt.Sprintf("%d", ms)
	}
	return fmt.Sprintf("%d.%d", ms, i%1000)
}
