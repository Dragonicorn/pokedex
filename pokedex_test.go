package main

import (
	"fmt"
	"testing"
	"time"

	pokecache "github.com/Dragonicorn/pokedex/internal"
)

// func TestCleanInput(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected []string
// 	}{
// 		{
// 			input:    "  hello  world  ",
// 			expected: []string{"hello", "world"},
// 		},
// 		{
// 			input:    "HeLp! I need somebody... ",
// 			expected: []string{"help!", "i", "need", "somebody..."},
// 		},
// 		{
// 			input:    "Charmander Bulbasaur PIKACHU",
// 			expected: []string{"charmander", "bulbasaur", "pikachu"},
// 		},
// 	}

// 	for _, c := range cases {
// 		actual := cleanInput(c.input)
// 		if len(actual) != len(c.expected) {
// 			t.Errorf("Length of actual and expected results array don't match.")
// 		}
// 		for i := range actual {
// 			word := actual[i]
// 			expectedWord := c.expected[i]
// 			//fmt.Printf("Actual Word: '%s', Expected Word: '%s'\n", word, expectedWord)
// 			if word != expectedWord {
// 				t.Errorf("Actual words and expected words don't match.")
// 			}
// 		}
// 	}
// }

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
