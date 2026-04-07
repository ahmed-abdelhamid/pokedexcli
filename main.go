// Package main implements the Pokedex CLI application.
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ahmed-abdelhamid/pokedexcli/internal/pokeapi"
)

func main() {
	cfg := &config{
		pokeapiClient: pokeapi.NewClient(5 * time.Minute),
		Pokedex:       make(map[string]pokeapi.Pokemon),
	}
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		cmd, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		args := words[1:]
		err := cmd.callback(cfg, args)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
