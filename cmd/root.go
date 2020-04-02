package cmd

import (
	"github.com/spf13/cobra"

	. "github.com/outerdev/algoc/config"
	. "github.com/outerdev/algoc/errors"
	"github.com/outerdev/algoc/kmd"
	. "github.com/outerdev/algoc/prompt"
)

type Config struct {
	Token string `yaml:"token" action:"prompt"`
	Host  string `yaml:"host" action:"prompt,url"`
}

var config Config

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	SetConfigFileName("algoc.yaml")
	err := LoadConfig(&config)
	if IsConfigNotPresent(err) {
		if err := PromptForValues(&config); err != nil {
			Fatal(err)
		}
		if err := WriteConfig(config); err != nil {
			Fatal(err)
		}
	} else if err != nil {
		Fatal(err)
	}
}

func init() {
	initConfig()

	if configDir, err := LocateConfigDir(); err == nil {
		kmd.Start(configDir)
	} else {
		panic(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// TODO: Add logging here
		// Fatal(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "algoc",
	Short: "Small command line utility to interact with the Algorand network",
	Long: `Command line utility to interact with the Algorand network.

Manage and issue assets with a few commands.`,
}
