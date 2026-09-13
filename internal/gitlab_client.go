package internal

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

func buildGitlabAPIBaseURL(hostname string) (string, error) {
	trimmedHost := strings.TrimSpace(hostname)
	if trimmedHost == "" {
		return "", fmt.Errorf("hostname cannot be empty")
	}

	if !strings.Contains(trimmedHost, "://") {
		trimmedHost = "https://" + trimmedHost
	}

	baseURL, err := url.Parse(trimmedHost)
	if err != nil {
		return "", fmt.Errorf("invalid hostname %q: %w", hostname, err)
	}
	if baseURL.Host == "" {
		return "", fmt.Errorf("invalid hostname %q", hostname)
	}

	baseURL.Path = path.Join(baseURL.Path, "api", "v4")
	baseURL.RawQuery = ""
	baseURL.Fragment = ""

	return baseURL.String(), nil
}

// NewGitlabClientFromCredentials creates a GitLab client from explicit credentials.
func NewGitlabClientFromCredentials(gitlabCredentials GitlabLogin) (*gitlab.Client, error) {
	apiBaseURL, err := buildGitlabAPIBaseURL(gitlabCredentials.Hostname)
	if err != nil {
		return nil, err
	}

	gitlabClient, err := gitlab.NewClient(gitlabCredentials.Token, gitlab.WithBaseURL(apiBaseURL))
	if err != nil {
		return nil, fmt.Errorf("creating GitLab client: %w", err)
	}

	return gitlabClient, nil
}

// NewGitlabClient creates a GitLab client from the stored credentials.
func NewGitlabClient() (*gitlab.Client, error) {
	gitlabCredentials, err := getGitlabCredentials()
	if err != nil {
		return nil, err
	}

	return NewGitlabClientFromCredentials(gitlabCredentials)
}
