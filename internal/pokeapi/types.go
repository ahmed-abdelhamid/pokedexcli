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

// Pokemon represents detailed information about a single pokemon.
type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

// PokemonStat represents a single stat entry for a pokemon.
type PokemonStat struct {
	BaseStat int     `json:"base_stat"`
	Stat     StatRef `json:"stat"`
}

// StatRef is a reference to a stat by name.
type StatRef struct {
	Name string `json:"name"`
}

// PokemonType represents a single type entry for a pokemon.
type PokemonType struct {
	Type TypeRef `json:"type"`
}

// TypeRef is a reference to a type by name.
type TypeRef struct {
	Name string `json:"name"`
}
