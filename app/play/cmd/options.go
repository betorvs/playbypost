/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// optionsCmd represents the options command
var optionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"opt", "o"},
	Short:   "get options from play by post server",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		channel := viper.GetString("channel")
		userid := viper.GetString("user-id")
		text := "opt"
		body, err := app.Web.PostCommand(userid, text, channel)
		if err != nil {
			app.Logger.Error("post command failed", "error", err.Error())
			os.Exit(1)
		}
		var msg types.Composed
		err = json.Unmarshal(body, &msg)
		if err != nil {
			app.Logger.Error("json unmarsharl error", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info("post command works", "answer", msg.Msg)
		if len(msg.Opts) > 0 {
			app.Logger.Info("Options", "options", msg.Opts)
			for _, v := range msg.Opts {
				fmt.Printf("Call: play exec '%s;%d'\n", v.Name, v.ID)
			}
		} else {
			app.Logger.Info("No options for you")
		}
	},
}

func init() {
	rootCmd.AddCommand(optionsCmd)
}
