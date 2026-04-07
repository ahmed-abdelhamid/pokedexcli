package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplit(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"leading and trailing spaces": {
			input: "  hello  world  ",
			want:  []string{"hello", "world"},
		},
		"mixed case words": {
			input: "Charmander Bulbasaur PIKACHU",
			want:  []string{"charmander", "bulbasaur", "pikachu"},
		},
		"empty string": {
			input: "",
			want:  []string{},
		},
		"only whitespace": {
			input: "   ",
			want:  []string{},
		},
		"single word": {
			input: "hello",
			want:  []string{"hello"},
		},
		"single word uppercase": {
			input: "HELLO",
			want:  []string{"hello"},
		},
		"tabs as whitespace": {
			input: "hello\tworld",
			want:  []string{"hello", "world"},
		},
		"multiple spaces between words": {
			input: "hello     world",
			want:  []string{"hello", "world"},
		},
		"newline separated": {
			input: "hello\nworld",
			want:  []string{"hello", "world"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf("%s", diff)
			}
		})
	}
}
