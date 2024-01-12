package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type GitlabLogin struct {
	Hostname string `json:"server"`
	Token    string `json:"token"`
}

var ConfigDir = filepath.Join(os.Getenv("HOME"), "/.config/ci-thief")

func getGitlabCredentials() (GitlabLogin, error) {
	gitlabCredentials := GitlabLogin{}
	file, err := os.ReadFile(filepath.Join(ConfigDir, "login.json"))
	if err != nil {
		return gitlabCredentials, fmt.Errorf("Could not find GitLab credentials file. Please run 'ci-thief login' first.")
	}

	if err := json.Unmarshal([]byte(file), &gitlabCredentials); err != nil {
		return gitlabCredentials, fmt.Errorf("Could not parse GitLab credentials file. Please run 'ci-thief login' again.")
	}

	return gitlabCredentials, nil
}
