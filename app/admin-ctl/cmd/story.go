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
			u, err := app.Web.GetStory()
			if err != nil {
				app.Logger.Error("get story", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range u {
				fmt.Printf("Story Title: %s, ID: %d, Master ID: %d, Announcement: \"%s\" \n", v.Title, v.ID, v.WriterID, v.Announcement)
			}

		case "create":
			body, err := app.Web.CreateStory(title, announcement, notes, writerID)
			if err != nil {
				app.Logger.Error("story error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// storyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// storyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	storyCmd.Flags().StringVarP(&title, "title", "t", "", " title to be used")
	storyCmd.Flags().StringVarP(&announcement, "announcement", "a", "", "announcement when starting story")
	storyCmd.Flags().StringVarP(&notes, "notes", "n", "-", " notes are internal notes about what it means and how to use it")
	storyCmd.Flags().IntVar(&writerID, "writer-id", 0, "writer id equal user ID")
}
