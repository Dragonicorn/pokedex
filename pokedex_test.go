package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HeLp! I need somebody... ",
			expected: []string{"help!", "i", "need", "somebody..."},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual and expected results array don't match.")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			//fmt.Printf("Actual Word: '%s', Expected Word: '%s'\n", word, expectedWord)
			if word != expectedWord {
				t.Errorf("Actual words and expected words don't match.")
			}
		}
	}
}
