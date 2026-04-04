package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type areaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	Next     string
	Previous string
}

// Print a goodbye message and exit the program
func commandExit(log *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// makeHelpCommand returns a function that prints the help message when called.
func makeHelpCommand(commands map[string]cliCommand) func(*config) error {
	return func(log *config) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Print("Usage:\n\n")

		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

// Fetch and print area data from the API, update the config with the next and previous URLs.
func commandMap(log *config) error {
	res, err := http.Get(log.Next)
	if err != nil {
		return fmt.Errorf("failed to fetch area data: %v", err)
	}
	defer res.Body.Close()

	// Create decoder, and decode the JSON response into a struct
	decoder := json.NewDecoder(res.Body)
	var areas areaResponse
	err = decoder.Decode(&areas)
	if err != nil {
		return fmt.Errorf("failed to decode area data: %v", err)
	}

	log.Next = areas.Next
	log.Previous = areas.Previous

	// Print the area names
	fmt.Println("Area Names:")
	for _, area := range areas.Results {
		fmt.Println("- ", area.Name)
	}

	return nil
}

// Fetch and print previous area data from the API, update the config with the next and previous URLs.
func commandMapBack(log *config) error {
	res, err := http.Get(log.Previous)
	if err != nil {
		return fmt.Errorf("failed to fetch previous area data: %v", err)
	}
	defer res.Body.Close()

	// Create decoder, and decode the JSON response into a struct
	decoder := json.NewDecoder(res.Body)
	var areas areaResponse
	err = decoder.Decode(&areas)
	if err != nil {
		return fmt.Errorf("failed to decode previous area data: %v", err)
	}

	log.Next = areas.Next
	log.Previous = areas.Previous

	// Print the area names
	fmt.Println("Previous Area Names:")
	for _, area := range areas.Results {
		fmt.Println("- ", area.Name)
	}

	return nil
}
