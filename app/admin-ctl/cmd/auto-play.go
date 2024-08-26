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

// autoPlayCmd represents the autoPlay command
var autoPlayCmd = &cobra.Command{
	Use:     "auto-play",
	Aliases: []string{"auto", "ap"},
	Short:   "creates or list a auto play from a story",
	Long: `

- create-by-title: will receive a encounter and next encounter title and create a auto play for you
	`,
	PreRun: loadApp,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "create":
			body, err := app.Web.CreateAutoPlay(displayText, storyID, solo)
			if err != nil {
				app.Logger.Error("autoPlay error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "next":
			body, err := app.Web.AddNextEncounter(autoPlayID, encounterID, nextEncounterID, displayText)
			if err != nil {
				app.Logger.Error("autoPlay next error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "next-by-title":
			encounters, err := app.getEncounters()
			if err != nil {
				os.Exit(1)
			}
			encounterID := app.findEncounterID(encounterTitle, encounters)
			if encounterID == 0 {
				app.Logger.Error("encounter not found", "title", encounterTitle)
				os.Exit(1)
			}
			nextEncounterID := app.findEncounterID(nextEncounterTitle, encounters)
			if nextEncounterID == 0 {
				app.Logger.Error("next encounter not found", "title", nextEncounterTitle)
				os.Exit(1)
			}
			body, err := app.Web.AddNextEncounter(autoPlayID, encounterID, nextEncounterID, displayText)
			if err != nil {
				app.Logger.Error("autoPlay next by title error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		}
	},
}

func init() {
	rootCmd.AddCommand(autoPlayCmd)
	autoPlayCmd.Flags().BoolVar(&solo, "solo", false, "solo adventure")
	autoPlayCmd.Flags().IntVar(&storyID, "story-id", 0, "story id")
	autoPlayCmd.Flags().StringVar(&displayText, "text", "", "text to be used")
	autoPlayCmd.Flags().IntVar(&autoPlayID, "auto-play-id", 0, "auto play id")
	autoPlayCmd.Flags().IntVar(&encounterID, "encounter-id", 0, "encounter id")
	autoPlayCmd.Flags().IntVar(&nextEncounterID, "next-encounter-id", 0, "next encounter id")
	autoPlayCmd.Flags().StringVar(&nextEncounterTitle, "next-encounter", "", "next encounter name")
	autoPlayCmd.Flags().StringVar(&encounterTitle, "encounter", "", "next encounter title")

}
