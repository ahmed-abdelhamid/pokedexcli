package main

import (
	"strings"

	"github.com/ahmed-abdelhamid/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient *pokeapi.Client
	Next          *string
	Previous      *string
	Pokedex       map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the Pokemon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(commands),
	}
	return commands
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}
