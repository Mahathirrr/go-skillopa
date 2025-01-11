package utils

import (
	"strings"
	"unicode"
)

func CreateSlug(title string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, title)
}

func SnakeCaseToTitle(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}