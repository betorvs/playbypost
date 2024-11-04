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

// storyCmd represents the story command
var storyCmd = &cobra.Command{
	Use:     "story [list|create]",
	Aliases: []string{"stories"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			u, err := app.getStories()
			if err != nil {
				os.Exit(1)
			}
			for _, v := range u {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Story", "id", v.ID, "title", v.Title, "writer_id", v.WriterID, "announcement", v.Announcement)

				}
			}

		case "create":
			body, err := app.Web.CreateStory(title, announcement, notes, writerID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("story error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "title", title, "master_id", writerID)

		case "update":
			body, err := app.Web.UpdateStory(title, announcement, notes, storyID, writerID)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("story update error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "title", title, "master_id", writerID)
		default:
			app.Logger.Info("story command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(storyCmd)
	storyCmd.Flags().StringVarP(&title, "title", "t", "", " title to be used")
	storyCmd.Flags().StringVarP(&announcement, "announcement", "a", "", "announcement when starting story")
	storyCmd.Flags().StringVarP(&notes, "notes", "n", "-", " notes are internal notes about what it means and how to use it")
	storyCmd.Flags().IntVar(&writerID, "writer-id", 0, "writer id equal user ID")
	storyCmd.Flags().IntVar(&storyID, "story-id", 0, "story id")
}
