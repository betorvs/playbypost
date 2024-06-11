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

// initiativeCmd represents the initiative command
var initiativeCmd = &cobra.Command{
	Use:     "initiative",
	Aliases: []string{"init", "i"},
	Short:   "generates a initiative",
	Long:    ``,
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		if encounterid == 0 {
			app.Logger.Error("encounter-id is mandatory")
			os.Exit(2)
		}
		if name == "" {
			app.Logger.Error("--name is mandatory")
			os.Exit(2)
		}
		body, err := app.Web.CreateInitiative(name, encounterid, isNPC)
		if err != nil {
			app.Logger.Error("initiative error", "error", err.Error())
			os.Exit(1)
		}
		var msg types.Msg
		err = json.Unmarshal(body, &msg)
		if err != nil {
			app.Logger.Error("initiative json unmarsharl error", "error", err.Error())
			os.Exit(1)
		}
		app.Logger.Info(msg.Msg, "encounter_id", encounterid)
	},
}

func init() {
	rootCmd.AddCommand(initiativeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initiativeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initiativeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initiativeCmd.Flags().StringVar(&name, "name", "", "initiative name")
	initiativeCmd.Flags().IntVar(&encounterid, "encounter-id", 0, "encounter ID")
	initiativeCmd.Flags().BoolVar(&isNPC, "non-player", false, "--non-player to add non player character as encounter participats")
}
