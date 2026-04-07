package main

import "fmt"

func commandMap(cfg *config) error {
	data, err := cfg.pokeapiClient.ListLocationAreas(cfg.Next)
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
