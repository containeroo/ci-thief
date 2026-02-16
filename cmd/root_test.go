package cmd

import "testing"

func TestShellQuote(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain",
			input: "hello",
			want:  "'hello'",
		},
		{
			name:  "contains single quote",
			input: "it's secret",
			want:  "'it'\\''s secret'",
		},
		{
			name:  "empty value",
			input: "",
			want:  "''",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := shellQuote(tc.input)
			if got != tc.want {
				t.Fatalf("shellQuote(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
