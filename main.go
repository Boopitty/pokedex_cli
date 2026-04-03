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
			err := cliCommands["exit"].callback()
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		case "help":
			err := cliCommands["help"].callback()
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}

		default:
			fmt.Printf("Unknown command\n")
		}
	}
}
