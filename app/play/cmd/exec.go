/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"x", "perform", "p"},
	Short:   "execute an action as a player in a channel",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		channel := viper.GetString("channel")
		userid := viper.GetString("user-id")
		tmpText := strings.Join(args, " ")
		var text string
		if strings.HasPrefix(tmpText, "choice") || strings.HasPrefix(tmpText, "cmd") {
			text = tmpText
		} else {
			app.Logger.Error("invalid command", "command", tmpText)
			os.Exit(1)
		}
		app.Logger.Info("text", "text", text)
		msg, err := app.Web.PostCommandComposed(userid, text, channel)
		if err != nil {
			app.Logger.Error("post command failed", "error", err.Error(), "text", text)
			os.Exit(1)
		}
		app.Logger.Info("post command works", "answer", msg.Msg)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
