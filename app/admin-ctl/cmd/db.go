/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
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
	Use:   "db [create|up|drop|verify]",
	Short: "Import DB Schema, create DB or Drop DB",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "create":
			create()
		case "up":
			up()
		case "down":
			down()
		case "drop":
			drop()
		case "ping":
			ping()
		case "verify":
			if len(args) < 2 {
				fmt.Println("error: table name is required")
				os.Exit(1)
			}
			verify(args[1])
		case "force":
			if len(args) < 2 {
				fmt.Println("error: version is required")
				os.Exit(1)
			}
			force(args[1])
		case "disconnect":
			disconnect()
		case "seed":
			password := utils.GetEnv("PGPASSWORD", "mypassword")
			if random {
				password = utils.RandomString(12)
			}
			{
				values := fmt.Sprintf("PGUSER=\"postgres\"\nPGPASSWORD=\"%s\"\nPGHOST=\"127.0.0.1\"\nPGPORT=5432\nPGDATABASE=\"playbypost\"", password)
				err := utils.Save(values, envFileTask)
				if err != nil {
					fmt.Println("error saving file", envFileTask, "error", err)
				}
				fmt.Println("file saved", envFileTask)
			}
			{
				values := fmt.Sprintf("export PGUSER=\"postgres\"\nexport PGPASSWORD=\"%s\"\nexport PGHOST=\"127.0.0.1\"\nexport PGPORT=5432\nexport PGDATABASE=\"playbypost\"", password)
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
	dbCmd.Flags().BoolVar(&random, "random", false, "Generate random password")
}

func ping() {
	creds := utils.LoadDBEnvVars()
	conn := strings.ReplaceAll(creds, "playbypost", "postgres")
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	err = db.Ping()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("ping executed ")
}

func create() {
	creds := utils.LoadDBEnvVars()
	conn := strings.ReplaceAll(creds, "playbypost", "postgres")
	fmt.Println(conn)
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
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
	defer func() {
		err := fs.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	dbpg, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := dbpg.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
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

func down() {
	creds := utils.LoadDBEnvVars()
	fs, err := data.Assets()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := fs.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	dbpg, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := dbpg.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
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
	err = m.Steps(-1)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
}

func drop() {
	creds := utils.LoadDBEnvVars()
	conn := strings.ReplaceAll(creds, "playbypost", "postgres")
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	_, err = db.ExecContext(ctx, "DROP DATABASE playbypost;")
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("drop executed ")
}

func verify(table string) {
	creds := utils.LoadDBEnvVars()
	db, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Verify columns
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '%s'", table))
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	fmt.Printf("Schema for table %s:\n", table)
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("  %s: %s\n", columnName, dataType)
	}

	// Verify foreign key constraints
	fkRows, err := db.QueryContext(ctx, fmt.Sprintf(`
		SELECT
			tc.constraint_name, kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.delete_rule
		FROM
			information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
			  ON tc.constraint_name = kcu.constraint_name
			  AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage AS ccu
			  ON ccu.constraint_name = tc.constraint_name
			  AND ccu.table_schema = tc.table_schema
			JOIN information_schema.referential_constraints AS rc
			  ON tc.constraint_name = rc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = '%s';
	`, table))
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := fkRows.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()

	fmt.Printf("\nForeign Key Constraints for table %s:\n", table)
	foundFK := false
	for fkRows.Next() {
		foundFK = true
		var constraintName, columnName, foreignTableName, foreignColumnName, deleteRule string
		if err := fkRows.Scan(&constraintName, &columnName, &foreignTableName, &foreignColumnName, &deleteRule); err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("  Constraint: %s, Column: %s, References: %s(%s), On Delete: %s\n", constraintName, columnName, foreignTableName, foreignColumnName, deleteRule)
	}
	if !foundFK {
		fmt.Println("  No foreign key constraints found.")
	}
}

func force(version string) {
	creds := utils.LoadDBEnvVars()
	fs, err := data.Assets()
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := fs.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	dbpg, err := pg.New(creds)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := dbpg.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
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
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	err = m.Force(versionInt)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
}

func disconnect() {
	creds := utils.LoadDBEnvVars()
	conn := strings.ReplaceAll(creds, "playbypost", "postgres")
	db, err := pg.New(conn)
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(2)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error ", err.Error())
			os.Exit(1)
		}
	}()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	_, err = db.ExecContext(ctx, "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'playbypost' AND pid <> pg_backend_pid();")
	if err != nil {
		fmt.Println("error ", err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnect executed ")
}
