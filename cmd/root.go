package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/containeroo/ci-thief/internal"
	"github.com/spf13/cobra"
	"gitlab.com/gitlab-org/api/client-go"
)

var (
	fetchNonRecursive bool
	gitlabClient      *gitlab.Client
)

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func appendOutput(output map[string][]string, scopeOrder *[]string, scope, key, value string) {
	if _, exists := output[scope]; !exists {
		*scopeOrder = append(*scopeOrder, scope)
	}

	output[scope] = append(output[scope], fmt.Sprintf("export %s=%s", key, shellQuote(value)))
}

func fetchProjectVars(projectID string, output map[string][]string, scopeOrder *[]string) error {
	opt := &gitlab.ListProjectVariablesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	for {
		projectVars, resp, err := gitlabClient.ProjectVariables.ListVariables(projectID, opt)
		if err != nil {
			return fmt.Errorf("could not fetch project variables: %w", err)
		}

		for _, v := range projectVars {
			appendOutput(output, scopeOrder, v.EnvironmentScope, v.Key, v.Value)
		}

		if resp == nil || resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil
}

func fetchGroupVars(groupID int64, output map[string][]string, scopeOrder *[]string) error {
	group, _, err := gitlabClient.Groups.GetGroup(groupID, nil)
	if err != nil {
		return fmt.Errorf("could not fetch group %d: %w", groupID, err)
	}

	if group.ParentID != 0 {
		if err := fetchGroupVars(group.ParentID, output, scopeOrder); err != nil {
			return err
		}
	}

	opt := &gitlab.ListGroupVariablesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	for {
		groupVars, resp, err := gitlabClient.GroupVariables.ListVariables(groupID, opt)
		if err != nil {
			return fmt.Errorf("could not fetch group variables: %w", err)
		}

		for _, v := range groupVars {
			fullScope := fmt.Sprintf("%s/%s", group.FullPath, v.EnvironmentScope)
			appendOutput(output, scopeOrder, fullScope, v.Key, v.Value)
		}

		if resp == nil || resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil
}

var RootCmd = &cobra.Command{
	Use:   "ci-thief [PROJECT_ID]",
	Short: "Fetch GitLab CI variables from a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		gitlabClient, err = internal.NewGitlabClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		output := make(map[string][]string)
		scopeOrder := make([]string, 0)

		if !fetchNonRecursive {
			project, _, err := gitlabClient.Projects.GetProject(args[0], nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if project.Namespace.Kind == "group" {
				if err := fetchGroupVars(project.Namespace.ID, output, &scopeOrder); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}

		if err := fetchProjectVars(args[0], output, &scopeOrder); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, scope := range scopeOrder {
			lines := output[scope]
			fmt.Printf("# env scope %s\n", scope)
			for _, line := range lines {
				fmt.Println(line)
			}
		}
	},
}

func Execute() {
	RootCmd.Flags().BoolVarP(&fetchNonRecursive, "non-recursive", "R", false, "Do not fetch variables from parent groups")
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
