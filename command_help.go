package main

import "fmt"

func commandHelp(commands map[string]cliCommand) func(*config) error {
	return func(_ *config) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		fmt.Println()
		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}
