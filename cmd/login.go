package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containeroo/ci-thief/internal"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitLab",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(internal.ConfigDir, 0o700); err != nil {
			fmt.Fprintln(os.Stderr, "Could not create config directory:", err)
			os.Exit(1)
		}

		gitlabToken := cmd.Flag("token").Value.String()
		if gitlabToken == "" {
			fmt.Print("Enter GitLab access token: ")
			gitlabTokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not read GitLab token:", err)
				os.Exit(1)
			}
			fmt.Println()

			gitlabToken = string(gitlabTokenBytes)
		}
		gitlabToken = strings.TrimSpace(gitlabToken)
		if gitlabToken == "" {
			fmt.Fprintln(os.Stderr, "GitLab token cannot be empty")
			os.Exit(1)
		}

		gitlabLogin := internal.GitlabLogin{
			Hostname: strings.TrimSpace(cmd.Flag("hostname").Value.String()),
			Token:    gitlabToken,
		}

		gitlabClient, err := internal.NewGitlabClientFromCredentials(gitlabLogin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not create GitLab client:", err)
			os.Exit(1)
		}
		currentUser, _, err := gitlabClient.Users.CurrentUser()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not login to GitLab:", err)
			os.Exit(1)
		}

		fileContent, err := json.Marshal(gitlabLogin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not marshal GitLab credentials:", err)
			os.Exit(1)
		}

		if err := os.WriteFile(filepath.Join(internal.ConfigDir, "login.json"), fileContent, 0o600); err != nil {
			fmt.Fprintln(os.Stderr, "Could not write GitLab credentials file:", err)
			os.Exit(1)
		}

		fmt.Printf("Login to %s successful as %s\n", gitlabLogin.Hostname, currentUser.Username)
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
