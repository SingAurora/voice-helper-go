package util

import "strings"

func RemoveInvisibleChars(input string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, input)
}
