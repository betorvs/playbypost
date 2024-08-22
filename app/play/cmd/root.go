/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string = "development"
	app     application
	server  string
	cfgFile string
)

type application struct {
	Logger *slog.Logger
	Web    *cli.Cli
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "play",
	Short:   "it's a command line interface to interact with play by post server as a player",
	Version: Version,
	Long: `
	Play by post is a server to play RPG games in a chat platform.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.play.yaml)")

	rootCmd.PersistentFlags().StringVar(&server, "server", "http://localhost:3000", "play by post http server, default: ")

}

func loadApp(cmd *cobra.Command, args []string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	h := cli.New(server)
	app.Logger = logger
	app.Web = h
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".play" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".play")
	}

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
