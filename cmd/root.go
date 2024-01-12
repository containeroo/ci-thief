package cmd

import (
	"fmt"
	"os"

	"github.com/containeroo/ci-thief/internal"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	fetchNonRecursive bool
	gitlabClient      *gitlab.Client
)

var output = make(map[string][]string)

func appendOutput(scope, key, value string) {
	output[scope] = append(output[scope], fmt.Sprintf("export %s='%s'", key, value))
}

func fetchProjectVars(projectID string) {
	projectVars, _, err := gitlabClient.ProjectVariables.ListVariables(projectID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, v := range projectVars {
		appendOutput(v.EnvironmentScope, v.Key, v.Value)
	}
}

func fetchGroupVars(groupID int) {
	group, _, err := gitlabClient.Groups.GetGroup(groupID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	groupVars, _, err := gitlabClient.GroupVariables.ListVariables(groupID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, v := range groupVars {
		fullScope := fmt.Sprintf("%s/%s", group.FullPath, v.EnvironmentScope)
		appendOutput(fullScope, v.Key, v.Value)
	}

	if group.ParentID != 0 {
		fetchGroupVars(group.ParentID)
	}
}

var RootCmd = &cobra.Command{
	Use:   "ci-thief [PROJECT_ID]",
	Short: "Fetch GitLab CI variables from a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		gitlabClient, err = internal.NewGitlabClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fetchProjectVars(args[0])

		if !fetchNonRecursive {
			project, _, err := gitlabClient.Projects.GetProject(args[0], nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if project.Namespace.Kind == "group" {
				fetchGroupVars(project.Namespace.ID)
			}
		}

		for k, v := range output {
			fmt.Printf("# env scope %s\n", k)
			for _, line := range v {
				fmt.Println(line)
			}
		}
	},
}

func Execute() {
	RootCmd.Flags().BoolVarP(&fetchNonRecursive, "non-recursive", "R", false, "Do not fetch variables from parent groups")
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
