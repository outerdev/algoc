package cmd

import (
	. "fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Algoc",
	Long:  `Prints the current version of the Algoc command line utility`,
	Run: func(cmd *cobra.Command, args []string) {
		Println("Version 0.1.0")
	},
}
