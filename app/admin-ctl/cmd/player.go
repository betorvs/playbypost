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
)

var ()

// playerCmd represents the player command
var playerCmd = &cobra.Command{
	Use:     "player [list|generate]",
	Aliases: []string{"players"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		switch args[0] {
		case "list":
			if storyid == 0 {
				app.Logger.Error("story-id is required")
				os.Exit(2)
			}
			u, err := app.Web.GetPlayersByStoryID(storyid)
			if err != nil {
				app.Logger.Error("get players", "error", err.Error())
				os.Exit(1)
			}
			for k, v := range u {
				fmt.Printf("PlayerID %d: %+v \n", k, v)
				fmt.Println("")
			}
		case "get":
			if playerid == 0 {
				app.Logger.Error("get player requires player id")
				os.Exit(1)
			}
			l, err := app.Web.GetPlayersByID(playerid)
			if err != nil {
				app.Logger.Error("get player by id", "error", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Player: %+v \n", l)
		case "generate":
			if playerid == 0 || storyid == 0 || name == "" {
				app.Logger.Error("all parameters are required")
				os.Exit(2)
			}
			body, err := app.Web.GeneratePlayer(name, playerid, storyid)
			if err != nil {
				app.Logger.Error("generate player error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "player", name)
		default:
			app.Logger.Info("players command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	playerCmd.Flags().IntVar(&playerid, "player-id", 0, "player id is equal user ID")
	playerCmd.Flags().IntVar(&storyid, "story-id", 0, "story ID")
	playerCmd.Flags().StringVar(&name, "name", "", "name of player, character's name")
}
