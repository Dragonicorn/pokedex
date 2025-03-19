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

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type Registry map[string]cliCommand

var registry Registry

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {
	registry = Registry{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

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
		cmd := words[0]
		command, ok := registry[cmd]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Println("Error returned by %s function.", command.name)
			}
		} else {
			fmt.Println("Unknown command.")
		}

	}
}
