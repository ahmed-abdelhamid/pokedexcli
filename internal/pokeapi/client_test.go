package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahmed-abdelhamid/pokedexcli/internal/pokecache"
	"github.com/google/go-cmp/cmp"
)

func TestListLocationAreas(t *testing.T) {
	t.Parallel()

	next := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
	validBody := `{
		"count": 1054,
		"next": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
		"previous": null,
		"results": [
			{"name": "canalave-city-area", "url": "https://pokeapi.co/api/v2/location-area/1/"},
			{"name": "eterna-city-area", "url": "https://pokeapi.co/api/v2/location-area/2/"}
		]
	}`

	tests := map[string]struct {
		pageURL    *string
		statusCode int
		body       string
		want       LocationAreaResponse
		wantErr    bool
	}{
		"first page (nil URL)": {
			pageURL:    nil,
			statusCode: http.StatusOK,
			body:       validBody,
			want: LocationAreaResponse{
				Count:    1054,
				Next:     &next,
				Previous: nil,
				Results: []LocationArea{
					{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
					{Name: "eterna-city-area", URL: "https://pokeapi.co/api/v2/location-area/2/"},
				},
			},
		},
		"explicit page URL": {
			pageURL:    &next,
			statusCode: http.StatusOK,
			body:       validBody,
			want: LocationAreaResponse{
				Count:    1054,
				Next:     &next,
				Previous: nil,
				Results: []LocationArea{
					{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
					{Name: "eterna-city-area", URL: "https://pokeapi.co/api/v2/location-area/2/"},
				},
			},
		},
		"server error": {
			pageURL:    nil,
			statusCode: http.StatusInternalServerError,
			body:       `not json`,
			wantErr:    true,
		},
		"invalid JSON": {
			pageURL:    nil,
			statusCode: http.StatusOK,
			body:       `{invalid`,
			wantErr:    true,
		},
		"empty results": {
			pageURL:    nil,
			statusCode: http.StatusOK,
			body:       `{"count":0,"next":null,"previous":null,"results":[]}`,
			want: LocationAreaResponse{
				Results: []LocationArea{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			cache := pokecache.NewCache(5 * time.Minute)
			defer cache.Stop()
			client := &Client{httpClient: srv.Client(), cache: cache}

			// When pageURL is nil the client builds the URL from baseURL,
			// which points at the real API. Override by providing the test
			// server URL so we always hit the fake.
			pageURL := tc.pageURL
			if pageURL == nil {
				u := srv.URL
				pageURL = &u
			} else {
				// Redirect the explicit URL to our test server too.
				u := srv.URL + "/location-area?offset=20&limit=20"
				pageURL = &u
			}

			got, err := client.ListLocationAreas(pageURL)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetLocationArea(t *testing.T) {
	t.Parallel()

	validBody := `{
		"pokemon_encounters": [
			{"pokemon": {"name": "tentacool", "url": "https://pokeapi.co/api/v2/pokemon/72/"}},
			{"pokemon": {"name": "magikarp", "url": "https://pokeapi.co/api/v2/pokemon/129/"}}
		]
	}`

	tests := map[string]struct {
		name string
		body string
		want LocationAreaDetail
	}{
		"valid area with pokemon": {
			name: "pastoria-city-area",
			body: validBody,
			want: LocationAreaDetail{
				PokemonEncounters: []PokemonEncounter{
					{Pokemon: PokemonRef{Name: "tentacool", URL: "https://pokeapi.co/api/v2/pokemon/72/"}},
					{Pokemon: PokemonRef{Name: "magikarp", URL: "https://pokeapi.co/api/v2/pokemon/129/"}},
				},
			},
		},
		"empty encounters": {
			name: "empty-area",
			body: `{"pokemon_encounters": []}`,
			want: LocationAreaDetail{
				PokemonEncounters: []PokemonEncounter{},
			},
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cache := pokecache.NewCache(5 * time.Minute)
			defer cache.Stop()
			client := &Client{httpClient: http.DefaultClient, cache: cache}

			// Seed cache so GetLocationArea returns without hitting the network.
			key := baseURL + "/location-area/" + tc.name
			client.cache.Add(key, []byte(tc.body))

			got, err := client.GetLocationArea(tc.name)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetPokemon(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		body string
		want Pokemon
	}{
		"pikachu": {
			name: "pikachu",
			body: `{
				"name": "pikachu",
				"base_experience": 112,
				"height": 4,
				"weight": 60,
				"stats": [{"base_stat": 35, "stat": {"name": "hp"}}],
				"types": [{"type": {"name": "electric"}}]
			}`,
			want: Pokemon{
				Name:           "pikachu",
				BaseExperience: 112,
				Height:         4,
				Weight:         60,
				Stats:          []PokemonStat{{BaseStat: 35, Stat: StatRef{Name: "hp"}}},
				Types:          []PokemonType{{Type: TypeRef{Name: "electric"}}},
			},
		},
		"mewtwo": {
			name: "mewtwo",
			body: `{
				"name": "mewtwo",
				"base_experience": 340,
				"height": 20,
				"weight": 1220,
				"stats": [{"base_stat": 106, "stat": {"name": "hp"}}],
				"types": [{"type": {"name": "psychic"}}]
			}`,
			want: Pokemon{
				Name:           "mewtwo",
				BaseExperience: 340,
				Height:         20,
				Weight:         1220,
				Stats:          []PokemonStat{{BaseStat: 106, Stat: StatRef{Name: "hp"}}},
				Types:          []PokemonType{{Type: TypeRef{Name: "psychic"}}},
			},
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cache := pokecache.NewCache(5 * time.Minute)
			defer cache.Stop()
			client := &Client{httpClient: http.DefaultClient, cache: cache}

			key := baseURL + "/pokemon/" + tc.name
			client.cache.Add(key, []byte(tc.body))

			got, err := client.GetPokemon(tc.name)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
