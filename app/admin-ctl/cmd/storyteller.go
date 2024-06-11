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

// storytellersCmd represents the storyteller command
var storytellersCmd = &cobra.Command{
	Use:     "storyteller [list|create]",
	Aliases: []string{"st", "writer"},
	Short:   "manipulate storytellers",
	Long:    `No long description yet`,
	Args:    cobra.ExactArgs(1),
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			u, err := app.Web.GetStoryteller()
			if err != nil {
				app.Logger.Error("get storyteller", "error", err.Error())
				os.Exit(1)
			}
			for _, v := range u {
				fmt.Printf("Username: %s , ID: %d \n", v.Username, v.ID)
			}

		case "create":
			body, err := app.Web.CreateStoryteller(username, userid, password)
			if err != nil {
				app.Logger.Error("create storytellersCmd", "error", err.Error())
				os.Exit(1)
			}
			var msg types.Msg
			err = json.Unmarshal(body, &msg)
			if err != nil {
				app.Logger.Error("json unmarsharl error", "error", err.Error())
				os.Exit(1)
			}
			app.Logger.Info(msg.Msg, "username", username, "userid", userid)
		default:
			app.Logger.Info("storytellers command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(storytellersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	storytellersCmd.Flags().StringVarP(&username, "username", "u", "", " username to be used")
	storytellersCmd.Flags().StringVar(&password, "password", "", "password should be unique")
}
