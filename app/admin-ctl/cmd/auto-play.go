/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
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
		case "list":
			autoPlays, err := app.Web.GetAutoPlay()
			if err != nil {
				app.Logger.Error("autoPlay list error", "error", err.Error())
				os.Exit(1)
			}
			for _, autoPlay := range autoPlays {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(autoPlay)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("autoPlay", "id", autoPlay.ID, "text", autoPlay.Text, "story_id", autoPlay.StoryID, "solo", autoPlay.Solo)

				}

			}

		case "create":
			body, err := app.Web.CreateAutoPlay(displayText, storyID, writerID, solo)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("autoPlay error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "next":
			// TODO: if dice roll objective, require update input with array of values, it fails in validation
			next := types.Next{
				UpstreamID:      autoPlayID,
				EncounterID:     encounterID,
				NextEncounterID: nextEncounterID,
				Text:            displayText,
			}
			if objectiveKind != "" && slices.Contains(types.Objectives(), objectiveKind) {
				next.Objective.Kind = objectiveKind
				if len(objectiveValues) > 0 {
					next.Objective.Values = objectiveValues
				}
			}
			body, err := app.Web.AddNextEncounter(next)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("autoPlay next error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "next-by-title":
			// TODO: if dice roll objective, require update input with array of values, it fails in validation
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
			next := types.Next{
				UpstreamID:      autoPlayID,
				EncounterID:     encounterID,
				NextEncounterID: nextEncounterID,
				Text:            displayText,
			}
			if objectiveKind != "" && slices.Contains(types.Objectives(), objectiveKind) {
				next.Objective.Kind = objectiveKind
				if len(objectiveValues) > 0 {
					next.Objective.Values = objectiveValues
				}
			}
			body, err := app.Web.AddNextEncounter(next)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("autoPlay next by title error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "publish":
			body, err := app.Web.PublishAutoPlay(autoPlayID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("autoPlay publish error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "id", autoPlayID)

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
	autoPlayCmd.Flags().StringVar(&objectiveKind, "objective-kind", "", fmt.Sprintf("objective kind: %v", types.Objectives()))
	autoPlayCmd.Flags().IntSliceVar(&objectiveValues, "objective-values", []int{}, "objective values")
	autoPlayCmd.Flags().IntVar(&writerID, "writer-id", 0, "master id equal user ID")
}
