package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	pokecache "github.com/Dragonicorn/pokedex/internal"
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
	callback    func(*pdConfig, []string) error
}

type Registry map[string]cliCommand

var registry Registry

type pdConfig struct {
	Next  string
	Prev  string
	Cache pokecache.Cache
}

type pdLocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pdLocationAreas struct {
	Count   int              `json:"count"`
	Next    string           `json:"next"`
	Prev    string           `json:"previous"`
	Results []pdLocationArea `json:"results"`
}

type pdEncounter struct {
	Pokemon pdLocationArea `json:"pokemon"`
}

type pdLocationAreaEncounters struct {
	Name       string        `json:"name"`
	Encounters []pdEncounter `json:"pokemon_encounters"`
}

func commandExit(cfg *pdConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *pdConfig, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func mapArea(cfg *pdConfig, url string) error {
	if len(url) == 0 {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	// Check for cached data
	body, ok := cfg.Cache.Get(url)
	if ok {
		fmt.Printf("Data at %s retrieved from cache...", url)
	} else {
		res, err := http.Get(url)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			//log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			return err
		}
		if err != nil {
			//log.Fatal(err)
			return err
		}
		cfg.Cache.Add(url, body)
		fmt.Printf("Data at %s added to cache...", url)
	}
	data := pdLocationAreas{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	// fmt.Printf("Count: %d\n", data.Count)
	// fmt.Printf("Next: '%s'\n", data.Next)
	cfg.Next = data.Next
	// fmt.Printf("Prev: '%s'\n", data.Prev)
	cfg.Prev = data.Prev
	for i, _ := range data.Results {
		//fmt.Printf("'%s': '%s'\n", data.Results[i].Name, data.Results[i].URL)
		fmt.Println(data.Results[i].Name)
	}
	return nil
}

func commandMap(cfg *pdConfig, args []string) error {
	//fmt.Printf("Previous: %s; Next: %s\n", cfg.Prev, cfg.Next)
	err := mapArea(cfg, cfg.Next)
	return err
}

func commandMapB(cfg *pdConfig, args []string) error {
	if cfg.Prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	//fmt.Printf("Previous: %s; Next: %s\n", cfg.Prev, cfg.Next)
	err := mapArea(cfg, cfg.Prev)
	return err
}

func commandExplore(cfg *pdConfig, args []string) error {
	if len(args) > 0 {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0])
		// Check for cached data
		body, ok := cfg.Cache.Get(url)
		if ok {
			fmt.Printf("Data at %s retrieved from cache...\n", url)
		} else {
			res, err := http.Get(url)
			if err != nil {
				//log.Fatal(err)
				return err
			}
			body, err = io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode > 299 {
				//log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
				return err
			}
			if err != nil {
				//log.Fatal(err)
				return err
			}
			cfg.Cache.Add(url, body)
			fmt.Printf("Data at %s added to cache...\n", url)
		}
		data := pdLocationAreaEncounters{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		fmt.Printf("Exploring %s...\nFound Pokemon:\n", data.Name)
		for i, _ := range data.Encounters {
			//fmt.Printf("'%s': '%s'\n", data.Results[i].Name, data.Results[i].URL)
			fmt.Printf(" - %s\n", data.Encounters[i].Pokemon.Name)
		}
	} else {
		return fmt.Errorf("No area provided to explore.")
	}
	return nil
}

func main() {
	registry = Registry{
		"explore": {
			name:        "explore",
			description: "Displays a list of Pokemon in a location area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Displays a list of location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous list of location areas in the Pokemon world",
			callback:    commandMapB,
		},
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

	pdCfg := &pdConfig{
		Next:  "",
		Prev:  "",
		Cache: *pokecache.NewCache(time.Second * 15),
	}

	var (
		text string
		args []string
	)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text = scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		words := cleanInput(text)
		cmd := words[0]
		if len(words) > 1 {
			args = words[1:]
		}
		command, ok := registry[cmd]
		if ok {
			err := command.callback(pdCfg, args)
			if err != nil {
				fmt.Printf("Error '%v' returned by %s function.\n", err, command.name)
			}
		} else {
			fmt.Println("Unknown command.")
		}
	}
}
