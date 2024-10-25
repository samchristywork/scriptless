package main

import (
	"unicode"
	"io/ioutil"
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

func loadAsset(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return string(data)
}
