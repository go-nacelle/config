package config

import "unicode"

func isExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}
