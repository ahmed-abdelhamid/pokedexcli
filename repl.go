package main

import (
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func() error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println()
			for _, cmd := range commands {
				fmt.Printf("%s: %s\n", cmd.name, cmd.description)
			}
			return nil
		},
	}
	return commands
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}
