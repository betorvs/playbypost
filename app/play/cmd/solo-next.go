/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:     "next",
	Aliases: []string{"nx", "n"},
	Short:   "get solo next info from play by post server",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		channel := viper.GetString("channel")
		userid := viper.GetString("user-id")
		text := "solo-next"
		msg, err := app.Web.PostCommandComposed(userid, text, channel)
		if err != nil {
			app.Logger.Error("post command failed", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info("post command works", "answer", msg.Msg)
		if len(msg.Opts) > 0 {
			app.Logger.Info("Options", "options", msg.Opts)
			for _, v := range msg.Opts {
				fmt.Printf("Call: play exec 'choice;%s;%d'\n", v.Name, v.ID)
			}
		} else {
			app.Logger.Info("No options for you")
		}
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
