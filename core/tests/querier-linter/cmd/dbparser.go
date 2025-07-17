/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/betorvs/playbypost/core/tests/querier-linter/core/finder"
	"github.com/spf13/cobra"
)

var (
	migrationFile, queryString, queryDir string
)

type Report struct {
	File    string
	Queries []string
	Errors  map[string]string
}

func NewReport(file string, queries []string) *Report {
	return &Report{
		File:    file,
		Queries: queries,
		Errors:  make(map[string]string),
	}
}

// dbparserCmd represents the dbparser command
var dbparserCmd = &cobra.Command{
	Use:     "dbparser",
	Aliases: []string{"dbp"},
	Short:   "this command will parse the db migration file and check queries commands",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		migration, err := finder.ParseAllDBMigrations("core/sys/db/data/migrations")
		if err != nil {
			fmt.Println("Error parsing Migration", err)
		}
		// fmt.Printf("Migration file parsed. Tables (%d) %v\n", len(migration.Tables), migration.Tables)
		// os.Exit(0)
		reports := []*Report{}
		// queries := []string{}
		if queryDir != "" {
			list, err := finder.ReadDir(queryDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, file := range list {
				codeFile := fmt.Sprintf("%v/%v", queryDir, file)
				// fmt.Println("codeFile: ", codeFile)
				res, err := finder.QueriesFromDir(codeFile)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				// queries = append(queries, res...)
				reports = append(reports, NewReport(codeFile, res))
			}
			// for _, q := range queries {
			// 	fmt.Println("query: ", q)
			// }
			// os.Exit(0)

		} else {
			// queries = append(queries, queryString)
			reports = append(reports, NewReport("query-string", []string{queryString}))
		}
		// errorsMap := make(map[string]string)
		for _, r := range reports {
			// fmt.Println("report: ", r.File)
			for _, q := range r.Queries {
				if q == "" {
					continue
				}
				tables, err := finder.ParseOne(q)
				if err != nil {
					fmt.Println("Error parsing query from source", err)
					r.Errors[r.File] = fmt.Sprintf("Error parsing query from source: %v, and query %s", err, q)
					// os.Exit(1)
				}
				for _, values := range tables.Tables {
					// fmt.Printf("columns in query: %v\n", values)
					for _, v := range values.Columns {
						_, ok := migration.Tables[values.Name]
						// fmt.Println("migration.Tables[values.Name]: ", migration.Tables[values.Name])
						if ok {
							lower := strings.ToLower(v)
							if !slices.Contains(migration.Tables[values.Name].Columns, lower) {
								// fmt.Printf("column %v not found in table %v\n", lower, values.Name)
								// errorsMap[values.Name] = fmt.Sprintf("column %v not found in table %v\n", lower, values.Name)
								r.Errors[values.Name] = fmt.Sprintf("column %v not found in table %v from query %s\n", lower, values.Name, q)
							}
						} else {
							// fmt.Printf("table %v not found in migration\n", values.Name)
							// errorsMap[values.Name] = fmt.Sprintf("table %v not found in migration\n", values.Name)
							r.Errors[values.Name] = fmt.Sprintf("table %v not found in migration from query %s\n", values.Name, q)
						}
					}
				}
			}

		}
		// fmt.Println("Errors found: ", errorsMap)
		for _, r := range reports {
			if len(r.Errors) != 0 {
				fmt.Println("Errors found:")
				for k, v := range r.Errors {
					fmt.Printf("File: %s Table: %v - Error: %v\n", r.File, k, v)
				}
				os.Exit(1)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(dbparserCmd)

	dbparserCmd.Flags().StringVar(&queryString, "query", "", "Query SQL string")
	dbparserCmd.Flags().StringVar(&queryDir, "dir", "", "Directory with go files with required annotations")
}
