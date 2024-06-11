/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/spf13/cobra"
)

var (
	Version string = "development"
	app     application
)

type application struct {
	Logger *slog.Logger
	Web    *cli.Cli
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "admin-ctl",
	Short:   "admin-ctl is used to control playbypost using http commands",
	Version: Version,
	Long: `Many commands will be easy to use in chat,
but we can have a easy way to manipulate it directly without waiting a bot to parse it. 

For example:

`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.admin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&server, "server", "http://localhost:3000", "playbypost http server, default: ")
	rootCmd.PersistentFlags().StringVar(&adminToken, "token", "", "admin token to access playbypost http server")
	rootCmd.PersistentFlags().StringVar(&adminUser, "admin-user", "admin", "admin token to access playbypost http server")

}

func loadApp(cmd *cobra.Command, args []string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	h := cli.NewHeaders(server, adminUser, adminToken)
	app.Logger = logger
	app.Web = h
}
