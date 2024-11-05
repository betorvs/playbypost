/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/betorvs/playbypost/core/utils"
	"github.com/spf13/cobra"
)

// stageCmd represents the stage command
var stageCmd = &cobra.Command{
	Use:     "stage [list|create|create-by-title]",
	Aliases: []string{"stages"},
	Short:   "A brief description of your command",
	Long: `

- create-by-title: will receive a story title and create a stage for this story
	`,
	PreRun: loadApp,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			t, err := app.getStages()
			if err != nil {
				os.Exit(1)
			}
			for _, v := range t {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Stage", "id", v.ID, "text", v.Text, "story_id", v.StoryID, "storyteller_id", v.StorytellerID)
				}

			}

		case "create":
			body, err := app.Web.CreateStage(displayText, userID, storyID, writerID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("stage error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "create-by-title":
			stories, err := app.getStories()
			if err != nil {
				os.Exit(1)
			}
			storyID = app.findStoryID(storyTitle, stories)
			if storyID == 0 {
				app.Logger.Error("story not found", "title", storyTitle)
				os.Exit(1)
			}
			body, err := app.Web.CreateStage(displayText, userID, storyID, writerID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("stage error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "add-encounter":
			app.Logger.Info("add-encounter command called")
			encounters, err := app.getEncounters()
			if err != nil {
				os.Exit(1)
			}
			encounterID := app.findEncounterID(encounterTitle, encounters)
			if encounterID == 0 {
				app.Logger.Error("encounter not found", "title", encounterTitle)
				os.Exit(1)
			}
			// stage
			stages, err := app.getStages()
			if err != nil {
				os.Exit(1)
			}
			var stageID int
			for _, v := range stages {
				if v.Text == stageTitle {
					stageID = v.ID
				}
			}
			if stageID == 0 {
				app.Logger.Error("stage not found", "title", stageTitle)
				os.Exit(1)
			}
			// story
			stories, err := app.getStories()
			if err != nil {
				os.Exit(1)
			}
			storyID = app.findStoryID(storyTitle, stories)
			if storyID == 0 {
				app.Logger.Error("story not found", "title", storyTitle)
				os.Exit(1)
			}
			body, err := app.Web.AddEncounterToStage(displayText, storyID, stageID, encounterID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("add encounter to stage error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		case "start":
			app.Logger.Info("start command called")
			// stage
			stages, err := app.getStages()
			if err != nil {
				os.Exit(1)
			}
			stageID := app.findStageID(stageTitle, stages)
			if stageID == 0 {
				app.Logger.Error("stage not found", "title", stageTitle)
				os.Exit(1)
			}
			body, err := app.Web.StartStage(stageID, chatChannelID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("start stage error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "text", displayText)

		default:
			app.Logger.Info("stage command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(stageCmd)
	stageCmd.Flags().StringVarP(&displayText, "display-text", "d", "", "display text for this stage")
	stageCmd.Flags().IntVar(&storyID, "story-id", 0, "story id from story created")
	stageCmd.Flags().StringVarP(&userID, "user-id", "u", "", "userid from chat integration")
	stageCmd.Flags().StringVarP(&storyTitle, "story-title", "t", "", "story title from story created")
	stageCmd.Flags().StringVarP(&stageTitle, "stage-title", "s", "", "stage title from stage created")
	stageCmd.Flags().StringVarP(&encounterTitle, "encounter-title", "e", "", "encounter title from encounter created")
	stageCmd.Flags().StringVarP(&chatChannelID, "channel-id", "c", "", "channel id from chat integration")
	stageCmd.Flags().IntVar(&writerID, "writer-id", 0, "master id equal user ID")
}
