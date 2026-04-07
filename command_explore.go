package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: explore <area_name>")
	}

	name := args[0]
	fmt.Printf("Exploring %s...\n", name)

	data, err := cfg.pokeapiClient.GetLocationArea(name)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}
