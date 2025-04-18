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

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:     "task [list|create]",
	Aliases: []string{"tasks"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  loadApp,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			t, err := app.Web.GetTask()
			if err != nil {
				app.Logger.Error("get task", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range t {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Task", "id", v.ID, "description", v.Description, "ability", v.Ability, "skill", v.Skill, "kind", v.Kind, "target", v.Target)

				}
			}

		case "create":
			body, err := app.Web.CreateTask(description, ability, skill, kind, target)
			if err != nil {
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("task error", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "description", description)
		default:
			app.Logger.Info("task command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().StringVarP(&description, "description", "d", "", " description to be used")
	taskCmd.Flags().StringVarP(&ability, "ability", "a", "", " ability to be used")
	taskCmd.Flags().StringVarP(&skill, "skill", "s", "", " skill to be used")
	taskCmd.Flags().IntVar(&kind, "kind", 2, "kind type: 2 for skill check and 3 for ability check")
	taskCmd.Flags().IntVar(&target, "target", 0, "target to be achieve for this task")
}
