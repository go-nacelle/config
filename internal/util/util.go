package util

import "unicode"

func IsExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}
