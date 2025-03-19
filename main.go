package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	// for i := range words {
	// 	fmt.Println(words[i])
	// }
	return words
}

func main() {
	var text string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text = scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		//fmt.Println(text)
		words := cleanInput(text)
		fmt.Printf("Your command was: %s\n", words[0])
	}
}
