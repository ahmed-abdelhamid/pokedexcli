package pokeapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

// Client wraps an HTTP client for interacting with the PokeAPI.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new PokeAPI client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// ListLocationAreas fetches a page of location areas from the given URL.
// If pageURL is nil, it fetches the first page.
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	var data LocationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaResponse{}, err
	}

	return data, nil
}
