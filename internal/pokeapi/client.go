package pokeapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ahmed-abdelhamid/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// Client wraps an HTTP client for interacting with the PokeAPI.
type Client struct {
	httpClient *http.Client
	cache      *pokecache.Cache
}

// NewClient returns a new PokeAPI client with a response cache.
// The cacheInterval controls how long responses are kept before expiration.
func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{},
		cache:      pokecache.NewCache(cacheInterval),
	}
}

// ListLocationAreas fetches a page of location areas from the given URL.
// If pageURL is nil, it fetches the first page. Responses are cached by URL.
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	body, err := c.fetch(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	var data LocationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaResponse{}, err
	}

	return data, nil
}

// GetLocationArea fetches details for a specific location area by name.
// Responses are cached by URL.
func (c *Client) GetLocationArea(name string) (LocationAreaDetail, error) {
	body, err := c.fetch(baseURL + "/location-area/" + name)
	if err != nil {
		return LocationAreaDetail{}, err
	}

	var data LocationAreaDetail
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaDetail{}, err
	}

	return data, nil
}

// GetPokemon fetches details for a pokemon by name. Responses are cached by URL.
func (c *Client) GetPokemon(name string) (Pokemon, error) {
	body, err := c.fetch(baseURL + "/pokemon/" + name)
	if err != nil {
		return Pokemon{}, err
	}

	var data Pokemon
	if err := json.Unmarshal(body, &data); err != nil {
		return Pokemon{}, err
	}

	return data, nil
}

// fetch retrieves the body for url, checking the cache first.
// On a cache miss it makes an HTTP GET and stores the response.
func (c *Client) fetch(url string) ([]byte, error) {
	if cached, ok := c.cache.Get(url); ok {
		return cached, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, body)

	return body, nil
}
