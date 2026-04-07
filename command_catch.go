package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: catch <pokemon_name>")
	}

	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	// Higher base experience → harder to catch.
	// Threshold is a random number in [0, 300).
	// If the threshold >= base experience, the catch succeeds.
	const maxExp = 300
	if rand.IntN(maxExp) < pokemon.BaseExperience { //nolint:gosec // game mechanic, not security
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	fmt.Printf("%s was caught!\n", name)
	fmt.Println("You may now inspect it with the inspect command.")
	cfg.Pokedex[name] = pokemon

	return nil
}
