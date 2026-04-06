package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

// areaResponse struct to match the JSON response from the API
type areaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Print a goodbye message and exit the program
func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// makeHelpCommand returns a function that prints the help message when called.
func makeHelpCommand(commands map[string]cliCommand) func(*config) error {
	return func(cfg *config) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Print("Usage:\n\n")

		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

// Fetch and print area data from the API, update the config with the next and previous URLs.
func commandMap(cfg *config) error {
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
func commandMapBack(cfg *config) error {
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
