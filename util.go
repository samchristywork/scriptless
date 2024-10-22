package main

import (
	"unicode"
)

func onlyLetters(input string) string {
	var result []rune
	for _, char := range input {
		if unicode.IsLetter(char) {
			result = append(result, char)
		}
	}
	return string(result)
}
