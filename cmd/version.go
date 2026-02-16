package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "2.3.2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ci-thief",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
