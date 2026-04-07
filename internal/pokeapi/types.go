// Package pokeapi provides a client for the PokeAPI.
package pokeapi

// LocationAreaResponse represents the paginated response from the location-area endpoint.
type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

// LocationArea represents a single location area entry.
type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// LocationAreaDetail represents the detailed response for a single location area.
type LocationAreaDetail struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter represents a pokemon found in a location area.
type PokemonEncounter struct {
	Pokemon PokemonRef `json:"pokemon"`
}

// PokemonRef is a reference to a pokemon by name and URL.
type PokemonRef struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
