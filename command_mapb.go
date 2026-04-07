package main

import "fmt"

func commandMapb(cfg *config, _ []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	data, err := cfg.pokeapiClient.ListLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	return nil
}
