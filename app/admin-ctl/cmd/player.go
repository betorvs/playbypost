/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
			u, err := app.Web.GetPlayers()
			if err != nil {
				app.Logger.Error("get players", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range u {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Player", "id", v.ID, "name", v.Name, "user_id", v.PlayerID, "stage_id", v.StageID)

				}
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
			b, _ := json.Marshal(l)
			fmt.Println(string(b))

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
			msg, err := app.Web.GeneratePlayer(name, playerUserID, 0, stageID)
			if err != nil {
				app.Logger.Error("generate player error", "error", err.Error(), "content", msg.Msg)
				os.Exit(1)
			}

			app.Logger.Info(msg.Msg, "player", name)

		case "generate":
			if playerID == 0 || storyID == 0 || name == "" {
				app.Logger.Error("all parameters are required")
				os.Exit(2)
			}
			msg, err := app.Web.GeneratePlayer(name, playerUserID, 0, storyID)
			if err != nil {
				app.Logger.Error("generate player error", "error", err.Error(), "content", msg.Msg)
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
	playerCmd.Flags().StringVar(&playerUserID, "player-id", "", "player id is equal user ID")
	playerCmd.Flags().IntVar(&stageID, "stage-id", 0, "stage ID")
	playerCmd.Flags().StringVar(&name, "name", "", "name of player, character's name")
	playerCmd.Flags().StringVar(&stageTitle, "stage-title", "", "title of stage")
}
