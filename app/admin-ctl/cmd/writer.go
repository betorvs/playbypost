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
				msg, _ := utils.ParseMsgBody(body)
				app.Logger.Error("create writersCmd", "error", err.Error(), "msg", msg.Msg)
				os.Exit(1)
			}
			msg, err := utils.ParseMsgBody(body)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "username", username)

		case "association":
			associations, err := app.Web.GetWriterUsersAssociation()
			if err != nil {
				app.Logger.Error("get writer users association", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range associations {
				app.Logger.Info("association", "id", v.ID, "writer_id", v.WriterID, "user_id", v.UserID)
			}
		case "delete-association":

			if id == 0 {
				app.Logger.Error("id is required")
				os.Exit(1)
			}

			err := app.Web.DeleteWriterUserAssociation(id)
			if err != nil {
				app.Logger.Error("delete writer user association", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info("writer user association deleted", "id", id)
		default:
			app.Logger.Info("writers command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(writersCmd)
	writersCmd.Flags().StringVarP(&username, "username", "u", "", " username to be used")
	writersCmd.Flags().StringVar(&password, "password", "", "password should be unique")
	writersCmd.Flags().IntVar(&id, "id", 0, "id of the writer user association to delete")
}
