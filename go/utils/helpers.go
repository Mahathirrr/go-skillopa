package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, " ")
}

