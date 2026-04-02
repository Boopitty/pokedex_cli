package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	// Annonyous struct to hold test cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "charmander Bulbasaur SQUIRTLE",
			expected: []string{"charmander", "bulbasaur", "squirtle"},
		},
		{
			input:    "Pikachu",
			expected: []string{"pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if len(actual) != len(c.expected) {
			t.Errorf("Expected length %d but got %d", len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected word '%s' but got '%s'", expectedWord, word)
			}
		}
	}
}
