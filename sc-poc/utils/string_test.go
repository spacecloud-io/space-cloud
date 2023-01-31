package utils

import (
	"testing"
)

func TestPluralize(t *testing.T) {

	tests := []struct {
		name string
		word string
		want string
	}{
		{
			name: "kind to resource",
			word: "CompiledGraphqlSource",
			want: "compiledgraphqlsources",
		},
		{
			name: "kind ending with 'y'",
			word: "OPAPolicy",
			want: "opapolicies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pluralize(tt.word)
			if got != tt.want {
				t.Errorf("Pluralize() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
