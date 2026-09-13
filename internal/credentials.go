// Package internal provides ci-thief's GitLab and credential helpers.
package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GitlabLogin contains the credentials used to access a GitLab instance.
type GitlabLogin struct {
	Hostname string `json:"server"`
	Token    string `json:"token"`
}

func defaultConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".config", "ci-thief")
	}

	return filepath.Join(homeDir, ".config", "ci-thief")
}

// ConfigDir is the directory containing ci-thief's configuration files.
var ConfigDir = defaultConfigDir()

func getGitlabCredentials() (GitlabLogin, error) {
	gitlabCredentials := GitlabLogin{}
	file, err := os.ReadFile(filepath.Join(ConfigDir, "login.json"))
	if err != nil {
		return gitlabCredentials, fmt.Errorf("could not find GitLab credentials file, please run 'ci-thief login' first")
	}

	if err := json.Unmarshal([]byte(file), &gitlabCredentials); err != nil {
		return gitlabCredentials, fmt.Errorf("could not parse GitLab credentials file, please run 'ci-thief login' again")
	}

	return gitlabCredentials, nil
}
