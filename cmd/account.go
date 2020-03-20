package cmd

import (
	. "fmt"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts on the algorand blockchain",
	Long: `Manage accounts on the algorand blockchain
allowing you to manage your assets from this command line utility`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an account on the algorand blockchain",
	Long: `Creates an account on the algorand blockchain
allowing you to manage your assets from this command line utility`,
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Println("Create account")
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(createCmd)
}
