/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/betorvs/playbypost/core/utils"
	"github.com/spf13/cobra"
)

// initiativeCmd represents the initiative command
var initiativeCmd = &cobra.Command{
	Use:     "initiative",
	Aliases: []string{"init", "i"},
	Short:   "generates a initiative",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		if encounterID == 0 {
			app.Logger.Error("encounter-id is mandatory")
			os.Exit(2)
		}
		if name == "" {
			app.Logger.Error("--name is mandatory")
			os.Exit(2)
		}
		body, err := app.Web.CreateInitiative(userID, channel, encounterID)
		if err != nil {
			msg, _ := utils.ParseMsgBody(body)
			app.Logger.Error("initiative error", "error", err.Error(), "msg", msg.Msg)
			os.Exit(1)
		}
		msg, err := utils.ParseMsgBody(body)
		if err != nil {
			app.Logger.Error("initiative json unmarsharl error", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info(msg.Msg, "encounter_id", encounterID)
	},
}

func init() {
	rootCmd.AddCommand(initiativeCmd)
	initiativeCmd.Flags().StringVarP(&userID, "user-id", "u", "", "userid from chat integration")
	initiativeCmd.Flags().StringVarP(&channel, "channel", "c", "", "channel from chat integration")
	initiativeCmd.Flags().IntVar(&encounterID, "encounter-id", 0, "encounter ID")
	initiativeCmd.Flags().BoolVar(&isNPC, "non-player", false, "--non-player to add non player character as encounter participats")
}
