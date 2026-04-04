package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

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
	}

	// Add the help command to the cliCommands map
	// after the map is defined to avoid a circular reference.
	cliCommands["help"] = cliCommand{
		name:        "help",
		description: "Displayes a help message",
		callback:    makeHelpCommand(cliCommands),
	}

	// Initialize the config struct with the initial API endpoint
	log := &config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
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
		default:
			fmt.Printf("Unknown command\n")
		}
	}
}
