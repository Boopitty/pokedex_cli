package main

import "strings"

func cleanInput(text string) []string {
	// Split the input into words slice and trim whitespace
	// Convert the input to lowercase
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return []string{}
	}
	for i, word := range words {
		words[i] = strings.TrimSpace(word)
	}
	return words
}
