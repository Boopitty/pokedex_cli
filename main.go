package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Boopitty/pokedex_cli/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
	pokedex  map[string]pokemonData
}

func main() {
	// Create a scanner to read user input from the command line
	scanner := bufio.NewScanner(os.Stdin)

	// Initialize the config struct with the initial API endpoint
	log := &config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
		cache:    pokecache.NewCache(60 * time.Second), // Create a ne with a 5-second reaping interval := pokecache.NewCache(5 * time.Second)
		pokedex:  make(map[string]pokemonData),
	}

	// Define the CLI commands and their descriptions
	var cliCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display area names",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous area names",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Explore an area to see which pokemon can be found there",
			callback:    explore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Catch a pokemon to see its base experience",
			callback:    catch,
		},
	}

	// Add the help command to the cliCommands map
	// after the map is defined to avoid a circular reference.
	cliCommands["help"] = cliCommand{
		name:        "help",
		description: "Displayes a help message",
		callback:    makeHelpCommand(cliCommands),
	}

	// REPL loop (Read-Eval-Print Loop)
	for {
		fmt.Print("Pokedex > ")

		// Read user input
		if !scanner.Scan() {
			break
		}

		// Store input in a variable and clean it
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		// If the cleaned input is empty,
		// continue to the next iteration of the loop
		if len(cleanedInput) == 0 {
			continue
		}

		switch cleanedInput[0] {
		case "exit":
			err := cliCommands["exit"].callback(log)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		case "help":
			err := cliCommands["help"].callback(log)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		case "map":
			err := cliCommands["map"].callback(log)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		case "mapb":
			err := cliCommands["mapb"].callback(log)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		case "explore":
			err := cliCommands["explore"].callback(log, cleanedInput[1:]...)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		case "catch":
			err := cliCommands["catch"].callback(log, cleanedInput[1:]...)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		default:
			fmt.Printf("Unknown command\n")
		}
	}
}

func cleanInput(text string) []string {
	// Split the input into words slice and trim whitespace
	// Convert the input to lowercase
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return []string{}
	}
	for i, word := range words {
		words[i] = strings.TrimSpace(word)
	}
	return words
}
