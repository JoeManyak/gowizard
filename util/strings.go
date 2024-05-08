package util

import (
	"fmt"
	"strings"
)

func PascalToSnakeCase(str string) string {
	snakeStr := ""
	for i, r := range []rune(str) {
		if i > 0 && i != len([]rune(str)) && r >= 'A' && r <= 'Z' {
			snakeStr += "_"
		}

		snakeStr += strings.ToLower(string(r))
	}

	return snakeStr
}

func MakePrivateName(str string) string {
	runeStr := []rune(str)
	return strings.ToLower(string(runeStr[0])) + string(runeStr[1:])
}

func MakePublicName(str string) string {
	runeStr := []rune(str)
	return strings.ToUpper(string(runeStr[0])) + string(runeStr[1:])
}

func MakeString(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}
