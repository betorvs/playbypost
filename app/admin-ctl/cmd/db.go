/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/betorvs/playbypost/core/sys/db/data"
	"github.com/betorvs/playbypost/core/sys/db/pg"
	"github.com/betorvs/playbypost/core/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db [create|up|drop]",
	Short: "Import DB Schema, create DB or Drop DB",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "create":
			create()
		case "up":
			up()
		case "drop":
			drop()
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func create() {
	creds := utils.LoadDBEnvVars()
	conn := strings.Replace(creds, "playbypost", "postgres", -1)
	fmt.Println(conn)
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	res, err := db.ExecContext(ctx, "CREATE DATABASE playbypost;")
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("created executed ", res)
}

func up() {
	creds := utils.LoadDBEnvVars()
	fs, err := data.Assets()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer fs.Close()
	dbpg, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer dbpg.Close()
	driver, err := postgres.WithInstance(dbpg, &postgres.Config{})
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	m, err := migrate.NewWithInstance("iofs", fs, "playbypost", driver)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	err = m.Up()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
}

func drop() {
	creds := utils.LoadDBEnvVars()
	conn := strings.Replace(creds, "playbypost", "postgres", -1)
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer db.Close()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	_, err = db.ExecContext(ctx, "DROP DATABASE playbypost;")
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("drop executed ")
}
