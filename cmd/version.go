package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of envp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
