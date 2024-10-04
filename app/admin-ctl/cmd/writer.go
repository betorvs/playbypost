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

// writersCmd represents the writer command
var writersCmd = &cobra.Command{
	Use:     "writer [list|create]",
	Aliases: []string{"st", "writer"},
	Short:   "manipulate writers",
	Long:    `No long description yet`,
	Args:    cobra.ExactArgs(1),
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			u, err := app.Web.GetWriter()
			if err != nil {
				app.Logger.Error("get writer", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range u {
				switch outputFormat {
				case formatJSON:
					b, _ := json.Marshal(v)
					fmt.Println(string(b))

				case formatLog:
					app.Logger.Info("Writer", "id", v.ID, "username", v.Username)

				}
			}

		case "create":
			body, err := app.Web.CreateWriter(username, password)
			if err != nil {
				app.Logger.Error("create writersCmd", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "username", username)
		default:
			app.Logger.Info("writers command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(writersCmd)
	writersCmd.Flags().StringVarP(&username, "username", "u", "", " username to be used")
	writersCmd.Flags().StringVar(&password, "password", "", "password should be unique")
}
