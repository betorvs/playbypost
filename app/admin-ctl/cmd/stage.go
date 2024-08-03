/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/spf13/cobra"
)

// stageCmd represents the stage command
var stageCmd = &cobra.Command{
	Use:     "stage [list|create]",
	Aliases: []string{"stages"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			t, err := app.Web.GetStage()
			if err != nil {
				app.Logger.Error("get stage", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range t {
				fmt.Printf("Stage Text: %s, Story ID: %v, Storyteller ID: %v \n", v.Text, v.StoryID, v.StorytellerID)
			}

		case "create":
			body, err := app.Web.CreateStage(displayText, userID, storyID)
			if err != nil {
				app.Logger.Error("stage error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	stageCmd.Flags().StringVarP(&displayText, "display-text", "d", "", "display text for this stage")
	stageCmd.Flags().IntVar(&storyID, "story-id", 0, "story id from story created")
	stageCmd.Flags().StringVarP(&userID, "user-id", "u", "", "userid from chat integration")
}
