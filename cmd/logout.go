package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containeroo/ci-thief/internal"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitLab",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(filepath.Join(internal.ConfigDir, "login.json")); os.IsNotExist(err) {
			fmt.Println("Not logged in")
			return
		}
		if err := os.Remove(filepath.Join(internal.ConfigDir, "login.json")); err != nil {
			fmt.Println("Could not delete GitLab credentials file:", err)
			os.Exit(1)
		}
		fmt.Println("Logout successful")
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
