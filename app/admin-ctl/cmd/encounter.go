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

// encounterCmd represents the encounter command
var encounterCmd = &cobra.Command{
	Use:     "encounter [list|get|create]",
	Aliases: []string{"enc", "e"},
	Short:   "creates or list a encounter from a story",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			l, err := app.Web.GetEncounters()
			if err != nil {
				app.Logger.Error("get encounter", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range l {
				fmt.Printf("Encounter Title: %s, ID: %d, Notes: \"%s\", Announcement: \"%s\" \n", v.Title, v.ID, v.Notes, v.Announcement)
			}

		case "get":
			if encounterID == 0 {
				app.Logger.Error("get encounter requires encounter id")
				os.Exit(1)
			}
			l, err := app.Web.GetEncounterByID(encounterID)
			if err != nil {
				app.Logger.Error("get encounter by id", "error", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Encounter Title: %s, ID: %d, Notes: \"%s\", Announcement: \"%s\" \n", l.Title, l.ID, l.Notes, l.Announcement)

		// case "start":
		// 	if encounterID == 0 {
		// 		app.Logger.Error("change phase requires encounter id")
		// 		os.Exit(1)
		// 	}

		// 	err := app.ChangePhase(encounterID, types.Started)
		// 	if err != nil {
		// 		app.Logger.Error("change encounter phase error", "error", err.Error())
		// 		os.Exit(1)
		// 	}
		// 	app.Logger.Info("encounter phase changed", "phase", types.Started)

		// case "finish":
		// 	if encounterID == 0 {
		// 		app.Logger.Error("change phase requires encounter id")
		// 		os.Exit(1)
		// 	}

		// 	err := app.ChangePhase(encounterID, types.Finished)
		// 	if err != nil {
		// 		app.Logger.Error("change encounter phase error", "error", err.Error())
		// 		os.Exit(1)
		// 	}
		// 	app.Logger.Info("encounter phase changed", "phase", types.Finished)

		// case "reset":
		// 	if encounterID == 0 {
		// 		app.Logger.Error("change phase requires encounter id")
		// 		os.Exit(1)
		// 	}

		// 	err := app.ChangePhase(encounterID, types.Waiting)
		// 	if err != nil {
		// 		app.Logger.Error("change encounter phase error", "error", err.Error())
		// 		os.Exit(1)
		// 	}
		// 	app.Logger.Info("encounter phase changed", "phase", types.Waiting)

		case "create":
			if storyid == 0 || title == "" || announcement == "" {
				app.Logger.Error("story-id, title and announcement are mandatory")
				os.Exit(2)
			}
			body, err := app.Web.CreateEncounter(title, announcement, notes, storyid, writerID, firstEncounter, lastEncounter)
			if err != nil {
				app.Logger.Error("encounter error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "title", title, "story_id", storyid)

		// case "add":
		// 	if encounterID == 0 {
		// 		app.Logger.Error("encounter-id are mandatory")
		// 		os.Exit(2)
		// 	}
		// 	if len(listPlayersID) == 0 {
		// 		app.Logger.Error("--players-id are mandatory")
		// 		os.Exit(2)
		// 	}
		// 	body, err := app.Web.AddParticipants(encounterID, isNPC, listPlayersID)
		// 	if err != nil {
		// 		app.Logger.Error("initiative error", "error", err.Error())
		// 		os.Exit(1)
		// 	}
		// 	var msg types.Msg
		// 	err = json.Unmarshal(body, &msg)
		// 	if err != nil {
		// 		app.Logger.Error("initiative json unmarsharl error", "error", err.Error())
		// 		os.Exit(1)
		// 	}
		// 	app.Logger.Info(msg.Msg, "encounter_id", encounterID, "is_npc", isNPC)

		default:
			app.Logger.Info("encounters command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(encounterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encounterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encounterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	encounterCmd.Flags().StringVarP(&title, "title", "t", "", " title to be used")
	encounterCmd.Flags().StringVarP(&announcement, "announcement", "a", "", "announcement when starting story")
	encounterCmd.Flags().StringVarP(&notes, "notes", "n", "-", " notes are internal notes about what it means and how to use it")
	encounterCmd.Flags().IntVar(&storyid, "story-id", 0, "story ID")
	encounterCmd.Flags().IntVar(&encounterID, "encounter-id", 0, "encounter ID")
	encounterCmd.Flags().IntSliceVar(&listPlayersID, "players-id", []int{}, "players by id, split by comma without space")
	encounterCmd.Flags().BoolVar(&firstEncounter, "first-encounter", false, "--first-encounter to suggest it as first encounter")
	encounterCmd.Flags().BoolVar(&lastEncounter, "last-encounter", false, "--last-encounter to suggest it as last encounter")
	encounterCmd.Flags().IntVar(&writerID, "writer-id", 0, "master id equal user ID")
}

func (app *application) ChangePhase(encounterID int, phase types.Phase) error {
	err := app.Web.ChangeEncounterPhase(encounterID, int(phase))
	if err != nil {
		return err
	}
	return nil
}
