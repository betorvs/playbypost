/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/spf13/cobra"
)

// validatorCmd represents the validator command
var validatorCmd = &cobra.Command{
	Use:     "validator [list|request]",
	Aliases: []string{"v"},
	Short:   "request an story, auto play or stage to be validate",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list", "l":
			obj, err := app.Web.GetValidator()
			if err != nil {
				app.Logger.Error("validator error", "error", err.Error())
				return
			}
			switch outputFormat {
			case formatJSON:
				// b, _ := json.Marshal(obj)
				fmt.Println(string(obj))

			case formatLog:
				app.Logger.Info("Validator", "obj", string(obj))

			}
		case "request", "req":
			obj := ""
			switch objectKind {
			case "auto-play", "autoplay":
				obj = "autoplay"
			case "stage":
				obj = "stage"
			default:
				obj = "story"
			}
			body, err := app.Web.ValidatorPut(obj, objectID)
			if err != nil {
				app.Logger.Error("validator error", "error", err.Error())
				return
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				return
			}
			app.Logger.Info("request", "kind", obj, "id", objectID, "status", msg.Msg)

		}
	},
}

func init() {
	rootCmd.AddCommand(validatorCmd)
	validatorCmd.PersistentFlags().StringVar(&objectKind, "kind", "", "Kind of object to be validated")
	validatorCmd.PersistentFlags().IntVar(&objectID, "id", 0, "ID of object to be validated")
}
