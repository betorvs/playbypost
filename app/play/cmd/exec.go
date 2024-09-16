/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
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
	Short:   "execute an action as a player in a channel",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		channel := viper.GetString("channel")
		userid := viper.GetString("user-id")
		tmpText := strings.Join(args, " ")
		var text string
		if strings.Contains(tmpText, ";") {
			text = fmt.Sprintf("cmd;%s", tmpText)
		}
		app.Logger.Info("text", "text", text)
		body, err := app.Web.PostCommand(userid, text, channel)
		if err != nil {
			app.Logger.Error("post command failed", "error", err.Error(), "text", text)
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
}
