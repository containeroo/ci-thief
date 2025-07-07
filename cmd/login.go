package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containeroo/ci-thief/internal"
	"github.com/spf13/cobra"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitLab",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(internal.ConfigDir); os.IsNotExist(err) {
			err := os.Mkdir(internal.ConfigDir, 0o755)
			if err != nil {
				fmt.Println("Could not create config directory:", err)
				return
			}
		}

		gitlabToken := cmd.Flag("token").Value.String()
		if gitlabToken == "" {
			fmt.Print("Enter GitLab access token: ")
			gitlabTokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("Could not read GitLab token:", err)
				return
			}
			fmt.Println()

			gitlabToken = string(gitlabTokenBytes)
		}

		gitlabLogin := internal.GitlabLogin{
			Hostname: cmd.Flag("hostname").Value.String(),
			Token:    gitlabToken,
		}

		fileContent, err := json.Marshal(gitlabLogin)
		if err != nil {
			fmt.Println("Could not marshal GitLab credentials:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(internal.ConfigDir, "login.json"), fileContent, 0o644); err != nil {
			fmt.Println("Could not write GitLab credentials file:", err)
			os.Exit(1)
		}

		gitlabClient, err := internal.NewGitlabClient()
		if err != nil {
			fmt.Println("Could not create GitLab client:", err)
		}
		_, _, err = gitlabClient.Users.ListUsers(&gitlab.ListUsersOptions{})
		if err != nil {
			fmt.Println("Could not login to GitLab:", err)
			return
		}

		fmt.Printf("Login to %s successful", gitlabLogin.Hostname)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("hostname", "", "GitLab Hostname (e.g. gitlab.example.com)")
	loginCmd.Flags().String("token", "", "Personal access token (api scope)")
	if err := loginCmd.MarkFlagRequired("hostname"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
