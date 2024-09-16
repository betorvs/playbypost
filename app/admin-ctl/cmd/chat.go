/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"os"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:     "chat",
	Aliases: []string{"c"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "create":
			// .AddChatInformation(userid, username, i.ChannelID, types.Discord)
			body, err := app.Web.AddChatInformation(chatUserID, chatUserName, chatChannelID, types.Discord)
			if err != nil {
				app.Logger.Error("chat error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", chatUserName)

		case "list":
			body, err := app.Web.GetChatInformation()
			if err != nil {
				app.Logger.Error("chat list error", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range body {
				app.Logger.Info("chat list", "username", v.Username, "userid", v.UserID, "channel", v.Channel)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chatCmd.Flags().StringVarP(&chatUserID, "userid", "u", "", "User ID")
	chatCmd.Flags().StringVarP(&chatUserName, "username", "n", "", "User Name")
	chatCmd.Flags().StringVarP(&chatChannelID, "channelid", "c", "", "Channel ID")

}
