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
			l, err := app.getEncounters()
			if err != nil {
				os.Exit(1)
			}
			for _, v := range l {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Encounter", "id", v.ID, "text", v.Title, "story_id", v.StoryID, "writer_id", v.WriterID)

				}
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

		case "create":
			if storyID == 0 || title == "" || announcement == "" {
				app.Logger.Error("story-id, title and announcement are mandatory")
				os.Exit(2)
			}
			body, err := app.Web.CreateEncounter(title, announcement, notes, storyID, writerID, firstEncounter, lastEncounter)
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
			app.Logger.Info(msg.Msg, "title", title, "story_id", storyID)

		case "update":
			if storyID == 0 || title == "" || announcement == "" || encounterID == 0 {
				app.Logger.Error("story-id, title, announcement and encounter-id are mandatory")
				os.Exit(2)
			}
			body, err := app.Web.UpdateEncounter(title, announcement, notes, encounterID, storyID, writerID, firstEncounter, lastEncounter)
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
			app.Logger.Info(msg.Msg, "title", title, "story_id", storyID)

		default:
			app.Logger.Info("encounters command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(encounterCmd)
	encounterCmd.Flags().StringVarP(&title, "title", "t", "", " title to be used")
	encounterCmd.Flags().StringVarP(&announcement, "announcement", "a", "", "announcement when starting story")
	encounterCmd.Flags().StringVarP(&notes, "notes", "n", "-", " notes are internal notes about what it means and how to use it")
	encounterCmd.Flags().IntVar(&storyID, "story-id", 0, "story ID")
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
