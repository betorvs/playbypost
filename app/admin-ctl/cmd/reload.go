/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	reloadServer string
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload player cache",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app.Web.UpdateURL(reloadServer)
		err := app.Web.PutReload()
		if err != nil {
			app.Logger.Error("reload player server failed", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info("reload player server called")
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	reloadCmd.Flags().StringVarP(&reloadServer, "reload-server", "r", "http://localhost:8090", " title to be used")
}
