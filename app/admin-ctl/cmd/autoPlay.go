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
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
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
			encounters, err := app.Web.GetEncounters()
			if err != nil {
				app.Logger.Error("get encounters", "error", err.Error())
				os.Exit(1)
			}
			var encounterID int
			var nextEncounterID int
			for _, v := range encounters {
				if v.Title == encounterTitle {
					encounterID = v.ID
				}
				if v.Title == nextEncounterTitle {
					nextEncounterID = v.ID
				}
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autoPlayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// autoPlayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	autoPlayCmd.Flags().BoolVar(&solo, "solo", false, "solo adventure")
	autoPlayCmd.Flags().IntVar(&storyID, "story-id", 0, "story id")
	autoPlayCmd.Flags().StringVar(&displayText, "text", "", "text to be used")
	autoPlayCmd.Flags().IntVar(&autoPlayID, "auto-play-id", 0, "auto play id")
	autoPlayCmd.Flags().IntVar(&encounterID, "encounter-id", 0, "encounter id")
	autoPlayCmd.Flags().IntVar(&nextEncounterID, "next-encounter-id", 0, "next encounter id")
	autoPlayCmd.Flags().StringVar(&nextEncounterTitle, "next-encounter", "", "next encounter name")
	autoPlayCmd.Flags().StringVar(&encounterTitle, "encounter", "", "next encounter title")

}
