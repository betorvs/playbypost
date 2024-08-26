/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
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

const (
	envFileTask = ".env.task"
	envFile     = ".env"
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
		case "ping":
			ping()
		case "seed":
			random := utils.RandomString(12)
			{
				values := fmt.Sprintf("PGUSER=\"postgres\"\nPGPASSWORD=\"%s\"\nPGHOST=\"127.0.0.1\"\nPGPORT=5432\nPGDATABASE=\"playbypost\"", random)
				err := utils.Save(values, envFileTask)
				if err != nil {
					fmt.Println("error saving file", envFileTask, "error", err)
				}
				fmt.Println("file saved", envFileTask)
			}
			{
				values := fmt.Sprintf("export PGUSER=\"postgres\"\nexport PGPASSWORD=\"%s\"\nexport PGHOST=\"127.0.0.1\"\nexport PGPORT=5432\nexport PGDATABASE=\"playbypost\"", random)
				err := utils.Save(values, envFile)
				if err != nil {
					fmt.Println("error saving file", envFile, "error", err)
				}
				fmt.Println("file saved", envFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}

func ping() {
	creds := utils.LoadDBEnvVars()
	db, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("ping executed ")
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
