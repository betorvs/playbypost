/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"x", "perform", "p"},
	Short:   "execute an action as a player in a story",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		story := viper.GetString("story")
		userid := viper.GetString("user-id")
		// player := viper.GetInt("player-id")
		text := strings.Join(args, " ")
		body, err := app.Web.PostCommand(userid, story, text)
		if err != nil {
			app.Logger.Error("post command failed", "error", err.Error())
			os.Exit(1)
		}
		var msg types.Msg
		err = json.Unmarshal(body, &msg)
		if err != nil {
			app.Logger.Error("json unmarsharl error", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info("post command works", "answer", msg.Msg)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
