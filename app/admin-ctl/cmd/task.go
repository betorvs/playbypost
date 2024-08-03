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
				fmt.Printf("Task Description: %s, ability: %s, skill: %s, kind: %v, target: %v \n", v.Description, v.Ability, v.Skill, v.Kind, v.Target)
			}

		case "create":
			body, err := app.Web.CreateTask(description, ability, skill, kind, target)
			if err != nil {
				app.Logger.Error("task error", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// taskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// taskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	taskCmd.Flags().StringVarP(&description, "description", "d", "", " description to be used")
	taskCmd.Flags().StringVarP(&ability, "ability", "a", "", " ability to be used")
	taskCmd.Flags().StringVarP(&skill, "skill", "s", "", " skill to be used")
	taskCmd.Flags().IntVar(&kind, "kind", 2, "kind type: 2 for skill check and 3 for ability check")
	taskCmd.Flags().IntVar(&target, "target", 0, "target to be achieve for this task")
}
