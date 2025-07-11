package internal

import (
	"fmt"

	"gitlab.com/gitlab-org/api/client-go"
)

func NewGitlabClient() (*gitlab.Client, error) {
	gitlabCredentials, err := getGitlabCredentials()
	if err != nil {
		return nil, err
	}

	gitlabClient, err := gitlab.NewClient(gitlabCredentials.Token, gitlab.WithBaseURL(fmt.Sprintf("https://%s/api/v4", gitlabCredentials.Hostname)))
	if err != nil {
		fmt.Printf("Error creating GitLab client: %s", err)
		return nil, err
	}
	return gitlabClient, nil
}
