package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

type pdStat struct {
	BaseStat int `json:"base_stat"`
	// Effort   int            `json:"effort"`
	Stat pdLocationArea `json:"stat"`
}

type pdType struct {
	// Slot int            `json:"slot"`
	Type pdLocationArea `json:"type"`
}

type pdPokemon struct {
	Name           string   `json:"name"`
	BaseExperience int      `json:"base_experience"`
	Height         int      `json:"height"`
	Weight         int      `json:"weight"`
	Stats          []pdStat `json:"stats"`
	Types          []pdType `json:"types"`
}

type Pokedex map[string]pdPokemon

var pokedex Pokedex

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

func (cfg *pdConfig) getURL(url string) ([]byte, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("No URL provided to function getURL")
	}
	// Check for cached data
	body, ok := cfg.Cache.Get(url)
	if ok {
		//fmt.Printf("Data at %s retrieved from cache...\n", url)
	} else {
		res, err := http.Get(url)
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			//log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			return nil, err
		}
		if err != nil {
			//log.Fatal(err)
			return nil, err
		}
		cfg.Cache.Add(url, body)
		//fmt.Printf("Data at %s added to cache...\n", url)
	}
	return body, nil
}

func (cfg *pdConfig) mapArea(url string) error {
	if len(url) == 0 {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	body, err := cfg.getURL(url)
	if err != nil {
		return err
	}
	data := pdLocationAreas{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	// fmt.Printf("Count: %d\n", data.Count)
	cfg.Next = data.Next
	cfg.Prev = data.Prev
	for i, _ := range data.Results {
		fmt.Println(data.Results[i].Name)
	}
	return nil
}

func commandMap(cfg *pdConfig, args []string) error {
	return cfg.mapArea(cfg.Next)
}

func commandMapB(cfg *pdConfig, args []string) error {
	if cfg.Prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return cfg.mapArea(cfg.Prev)
}

func commandExplore(cfg *pdConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No area provided to explore.")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", args[0])
	body, err := cfg.getURL(url)
	if err != nil {
		return err
	}
	data := pdLocationAreaEncounters{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\nFound Pokemon:\n", data.Name)
	for i, _ := range data.Encounters {
		fmt.Printf(" - %s\n", data.Encounters[i].Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *pdConfig, args []string) error {
	//fmt.Printf("commandCatch: args: '%d'\n", len(args))
	if len(args) == 0 {
		return fmt.Errorf("No Pokemon name provided to catch.")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0])
	body, err := cfg.getURL(url)
	if err != nil {
		return err
	}
	data := pdPokemon{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)
	rn := rand.Uint32() >> 24
	if rn > uint32(data.BaseExperience) {
		fmt.Printf("%s was caught!\n", data.Name)
		pokedex[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
		if _, ok := pokedex[data.Name]; ok {
			delete(pokedex, data.Name)
		}
	}
	return nil
}

func commandInspect(cfg *pdConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No Pokemon name provided to inspect.")
	}
	if data, ok := pokedex[args[0]]; !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", data.Name)
		fmt.Printf("Height: %d\n", data.Height)
		fmt.Printf("Weight: %d\n", data.Weight)
		fmt.Println("Stats:")
		for _, s := range data.Stats {
			fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range data.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
	}
	//fmt.Printf("Pokedex: %v\n", pokedex)
	return nil
}

func main() {
	registry = Registry{
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    commandInspect,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon and add it to the Pokedex",
			callback:    commandCatch,
		},
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
	pokedex = make(map[string]pdPokemon, 0)

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
