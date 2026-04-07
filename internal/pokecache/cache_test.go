package pokecache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCache(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		key   string
		val   []byte
		found bool
	}{
		"add and get": {
			key:   "https://pokeapi.co/api/v2/location-area",
			val:   []byte(`{"results":[]}`),
			found: true,
		},
		"missing key": {
			key:   "",
			val:   nil,
			found: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := NewCache(5 * time.Minute)
			defer c.Stop()

			if tc.found {
				c.Add(tc.key, tc.val)
			}

			got, ok := c.Get(tc.key)
			if ok != tc.found {
				t.Fatalf("Get(%q) found = %v, want %v", tc.key, ok, tc.found)
			}
			if diff := cmp.Diff(tc.val, got); diff != "" {
				t.Fatalf("value mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCacheReap(t *testing.T) {
	t.Parallel()

	interval := 50 * time.Millisecond
	c := NewCache(interval)
	defer c.Stop()

	c.Add("key", []byte("value"))

	// Entry should exist immediately.
	if _, ok := c.Get("key"); !ok {
		t.Fatal("expected entry to exist before reap")
	}

	// Wait for at least two intervals so the reaper runs.
	time.Sleep(3 * interval)

	if _, ok := c.Get("key"); ok {
		t.Fatal("expected entry to be reaped")
	}
}
