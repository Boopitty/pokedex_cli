package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

// REPL (Read-Eval-Print Loop) for the Pokedex CLI application.
// This File contains the functions that the main REPL loop will call when the user enters a command.

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

// areaResponse struct to match the JSON response from the API
type areaResponse struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Count    int    `json:"count"`
}

// Partial struct to match needed info from the "explore" API response
type areaPokemon struct {
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemonData struct {
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`

	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`

	Name    string `json:"name"`
	URL     string `json:"url"`
	BaseExp int    `json:"base_experience"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
}

// commandExit prints a goodbye message and exits the program.
func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// makeHelpCommand returns a function that prints the help message when called.
func makeHelpCommand(commands map[string]cliCommand) func(*config, ...string) error {
	return func(cfg *config, args ...string) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Print("Usage:\n\n")

		for _, cmd := range commands {
			fmt.Printf("%s:\n - %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

// Fetch and print area data from the API, update the config with the next and previous URLs.
func commandMap(cfg *config, args ...string) error {
	// Get url from config and check if it's in the cache
	url := cfg.Next
	data, ok := cfg.cache.Get(url)
	var areas areaResponse

	// fetch and decode from the API if not in cache
	if !ok {

		fmt.Println("Fetching area data from the API...")
		// Make Request to the API
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch area data: %v", err)
		}
		defer res.Body.Close()

		// Read the response as a byte slice
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read area data: %v", err)
		}
	} else {
		fmt.Println("Area data found in cache, using cached data...")
	}

	// Decode the JSON response into the areas struct
	err := json.Unmarshal(data, &areas)
	if err != nil {
		return fmt.Errorf("failed to decode area data: %v", err)
	}

	// Update the config with the next and previous URLs
	cfg.Next = areas.Next
	cfg.Previous = areas.Previous

	// Add the data to the cache
	cfg.cache.Add(url, data)

	// Print the area names
	fmt.Println("Area Names:")
	for _, area := range areas.Results {
		fmt.Println("- ", area.Name)
	}

	return nil
}

// Fetch and print previous area data from the API, update the config with the next and previous URLs.
func commandMapBack(cfg *config, args ...string) error {
	url := cfg.Previous
	data, ok := cfg.cache.Get(url)
	var areas areaResponse

	// if url not in cache, fetch and read from the API
	if !ok {
		fmt.Println("Fetching previous area data from the API...")

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch previous area data: %v", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read previous area data: %v", err)
		}
	} else {
		fmt.Println("Previous area data found in cache, using cached data...")
	}

	err := json.Unmarshal(data, &areas)
	if err != nil {
		return fmt.Errorf("failed to decode previous area data: %v", err)
	}

	// Update the config
	cfg.Next = areas.Next
	cfg.Previous = areas.Previous
	cfg.cache.Add(url, data)

	// Print the area names
	fmt.Println("Previous Area Names:")
	for _, area := range areas.Results {
		fmt.Println("- ", area.Name)
	}

	return nil
}

// Takes a the name of an area and prints out the pokemon that can be found in that area.
func explore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no area name provided")
	}
	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName
	data, ok := cfg.cache.Get(url)

	// if url not in cache, fetch and read from the API
	if !ok {
		fmt.Printf("Fetching data for area '%s' from the API...\n", areaName)

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch area data: %v", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read area data: %v", err)
		}
	} else {
		fmt.Printf("Data for area '%s' found in cache, using cached data...\n", areaName)
	}

	var areaData areaPokemon
	err := json.Unmarshal(data, &areaData)
	if err != nil {
		return fmt.Errorf("failed to decode area data: %v", err)
	}

	cfg.cache.Add(url, data)

	fmt.Printf("Exploring '%s'...\n", areaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range areaData.Pokemon_encounters {
		fmt.Println("- ", encounter.Pokemon.Name)
	}
	return nil
}

// Takes a pokemon's name and attempts to catch it.
func catch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon name provided")
	}
	// Get proper URL and get response
	pokemon := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	data, ok := cfg.cache.Get(url)

	if !ok {
		fmt.Printf("Fetching data for pokemon '%s' from the API...\n", pokemon)
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch pokemon data: %v", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read pokemon data: %v", err)
		}
	} else {
		fmt.Printf("Data for pokemon '%s' found in cache, using cached data...\n", pokemon)
	}

	var pokemonData pokemonData
	err := json.Unmarshal(data, &pokemonData)
	if err != nil {
		return fmt.Errorf("failed to decode pokemon data: %v", err)
	}

	cfg.cache.Add(url, data)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonData.Name)

	chance := pokemonData.BaseExp * rand.Intn(100) / 100
	if chance < 50 {
		fmt.Printf("Congratulations! You caught %s!\n", pokemonData.Name)
		cfg.pokedex[pokemonData.Name] = pokemonData
	} else {
		fmt.Printf("Oh no! %s escaped!\n", pokemonData.Name)
	}

	return nil
}

func inspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon name provided")
	}

	// Check if the pokemon is in the pokedex
	name := args[0]
	pokemon, ok := cfg.pokedex[name]
	if !ok {
		return fmt.Errorf("you haven't caught %s yet!", name)
	}

	// Print the pokemon's details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExp)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("	- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("	- %s\n", t.Type.Name)
	}
	return nil
}

func pokedex(cfg *config, args ...string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your pokedex is empty! Catch some pokemon to see them here.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.pokedex {
		fmt.Printf("	- %s\n", name)
	}
	return nil
}
