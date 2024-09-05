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
	Use:     "player [list|generate|add-player-stage-by-title]",
	Aliases: []string{"players"},
	Short:   "A brief description of your command",
	Long: `
	
- add-player-stage-by-title: will receive a stage title and create a player for this stage
	`,
	PreRun: loadApp,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		switch args[0] {
		case "list":
			// if stageID == 0 {
			// 	app.Logger.Error("stage-id is required")
			// 	os.Exit(2)
			// }
			u, err := app.Web.GetPlayers()
			if err != nil {
				app.Logger.Error("get players", "error", err.Error())
				os.Exit(1)
			}
			for k, v := range u {
				fmt.Printf("PlayerID %d: %+v \n", k, v)
				fmt.Println("")
			}
		case "get":
			if playerID == 0 {
				app.Logger.Error("get player requires player id")
				os.Exit(1)
			}
			l, err := app.Web.GetPlayersByID(playerID)
			if err != nil {
				app.Logger.Error("get player by id", "error", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Player: %+v \n", l)
		case "add-player-stage-by-title":
			if stageTitle == "" {
				app.Logger.Error("stage-title is required")
				os.Exit(2)
			}
			stages, err := app.Web.GetStage()
			if err != nil {
				app.Logger.Error("get stages", "error", err.Error())
				os.Exit(1)
			}
			var stageID int
			for _, v := range stages {
				if v.Text == stageTitle {
					stageID = v.ID
				}
			}
			if stageID == 0 {
				app.Logger.Error("stage not found", "stage-title", stageTitle)
				os.Exit(1)
			}
			body, err := app.Web.GeneratePlayer(name, "", playerID, storyID)
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

		case "generate":
			if playerID == 0 || storyID == 0 || name == "" {
				app.Logger.Error("all parameters are required")
				os.Exit(2)
			}
			body, err := app.Web.GeneratePlayer(name, "", playerID, storyID)
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
	playerCmd.Flags().IntVar(&playerID, "player-id", 0, "player id is equal user ID")
	playerCmd.Flags().IntVar(&stageID, "stage-id", 0, "stage ID")
	playerCmd.Flags().StringVar(&name, "name", "", "name of player, character's name")
	playerCmd.Flags().StringVar(&stageTitle, "stage-title", "", "title of stage")
}
