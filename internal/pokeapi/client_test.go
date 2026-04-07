package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

			client := &Client{httpClient: srv.Client()}

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
