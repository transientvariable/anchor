package support

import (
	"strings"
	"unicode"
)

// RemoveWhitespace removes all whitespace characters from the provided string.
func RemoveWhitespace(s string) string {
	if len(s) == 0 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	for _, c := range s {
		if !unicode.IsSpace(c) {
			b.WriteRune(c)
		}
	}
	return b.String()
}
