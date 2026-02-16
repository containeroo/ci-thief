package internal

import "testing"

func TestBuildGitlabAPIBaseURL(t *testing.T) {
	testCases := []struct {
		name      string
		hostname  string
		want      string
		shouldErr bool
	}{
		{
			name:     "hostname only",
			hostname: "gitlab.com",
			want:     "https://gitlab.com/api/v4",
		},
		{
			name:     "full URL",
			hostname: "https://gitlab.example.com",
			want:     "https://gitlab.example.com/api/v4",
		},
		{
			name:     "full URL with path",
			hostname: "https://gitlab.example.com/gitlab",
			want:     "https://gitlab.example.com/gitlab/api/v4",
		},
		{
			name:      "invalid hostname",
			hostname:  "://broken",
			shouldErr: true,
		},
		{
			name:      "empty hostname",
			hostname:  "  ",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildGitlabAPIBaseURL(tc.hostname)
			if tc.shouldErr {
				if err == nil {
					t.Fatalf("buildGitlabAPIBaseURL(%q) expected error, got nil", tc.hostname)
				}
				return
			}

			if err != nil {
				t.Fatalf("buildGitlabAPIBaseURL(%q) unexpected error: %v", tc.hostname, err)
			}
			if got != tc.want {
				t.Fatalf("buildGitlabAPIBaseURL(%q) = %q, want %q", tc.hostname, got, tc.want)
			}
		})
	}
}
